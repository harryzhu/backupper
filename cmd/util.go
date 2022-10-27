package cmd

import (
	"sqlconf"
	"strconv"

	//"fmt"
	//"io"
	//"io/ioutil"
	//"net/http"
	"net/url"
	//"os"
	"path/filepath"

	//"strconv"
	"strings"
	"sync"

	//"time"

	//	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

var (
	IsOverwrite     bool
	BatchSize       int = 3
	DownloadSaveDir string
	URLFileList     map[string]string = make(map[string]string, 10)
)

func init() {

}

func urlToSavePath(URL string) (savePath string) {
	u, err := url.Parse(URL)
	if err != nil {
		logger.Error("url parse error:", zap.Error(err))
	}

	u_host := u.Host
	if strings.Contains(u.Host, ":") {
		u_host = strings.Split(u.Host, ":")[0]
	}

	DownloadSaveDir = strings.TrimRight(DownloadSaveDir, "/")
	uPath := strings.ReplaceAll(u.Path, "\\", "/")
	uPath2 := sqlconf.Filepathify(uPath)

	if DownloadSaveDir == "" {
		logger.Fatal("DownloadSaveDir cannot be empty")
	}

	savePath = filepath.Join(DownloadSaveDir, u_host, uPath2)
	saveDir := filepath.Dir(savePath)

	sqlconf.MakeDirs(saveDir)

	//logger.Info("urlToSavePath:", zap.String("url", URL), zap.String("savePath", savePath))
	return savePath
}

func prepareURLFileList(URL string) error {
	lst, err := sqlconf.GetURLContent(URL)
	if err != nil {
		logger.Error("prepareURLFileList", zap.Error(err))
		return err
	}

	logger.Info("prepareURLFileList", zap.String("url", URL))
	lst = strings.ReplaceAll(lst, "\r\n", "\n")
	arrList := strings.Split(lst, "\n")
	//fmt.Println(arrList)
	for _, line := range arrList {
		line = strings.Trim(line, " ")
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r")

		if strings.HasPrefix(line, "http") != true {
			//logger.Info("prepareURLFileList", zap.String("skip", line))
			continue
		}

		URLFileList[line] = urlToSavePath(line)
	}

	return nil
}

func runDownload() error {
	if len(URLFileList) == 0 {
		logger.Info("SKIP download", zap.String("StartDownload", "URLFileList is empty"))
		return nil
	}
	var totalTasks int = 0
	for u, f := range URLFileList {
		if u == "" || f == "" {
			continue
		}
		totalTasks += 1
		p := strings.Join([]string{"downloading:[", strconv.Itoa(totalTasks), "/", strconv.Itoa(len(URLFileList)), "]"}, "")

		sqlconf.DownloadFile(u, f, IsOverwrite, p)
	}

	return nil
}

func runMultiDownload() error {
	if len(URLFileList) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	var runningQueue int = 0
	var totalTasks int = 0
	var taskTitle string = ""

	if len(URLFileList) <= BatchSize {
		BatchSize = len(URLFileList)
	}
	logger.Info("StartMultiDownload", zap.Int("batch-size", BatchSize), zap.Int("task-count", len(URLFileList)))

	for u, f := range URLFileList {
		wg.Add(1)
		totalTasks += 1
		taskTitle = strings.Join([]string{"downloading:[", strconv.Itoa(totalTasks), "/", strconv.Itoa(len(URLFileList)), "]"}, "")
		//logger.Info(p, zap.String("url", u), zap.String("localPath", f))

		go func(u, f, taskTitle string) {
			sqlconf.DownloadFile(u, f, IsOverwrite, taskTitle)

			wg.Done()
		}(u, f, taskTitle)

		runningQueue++
		if runningQueue >= BatchSize {
			wg.Wait()
			runningQueue = 0
		}

	}
	wg.Wait()

	return nil
}
