package cmd

import (
	//"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	//	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

var (
	IsOverwrite bool = true
	BatchSize   int  = 3
	DirSaveRoot string
	URLFileList map[string]string = make(map[string]string, 10)
)

func init() {
	DirSaveRoot = strings.TrimRight(DirSaveRoot, "/")
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

	uPath := strings.ReplaceAll(u.Path, "\\", "/")
	uPath2 := Filepathify(uPath)

	savePath = filepath.Join(DirSaveRoot, u_host, uPath2)
	saveDir := filepath.Dir(savePath)

	makeDirs(saveDir)

	logger.Info("urlToSavePath:", zap.String("url", URL), zap.String("original-path", uPath), zap.String("filepathify-path", uPath2), zap.String("saveDir", saveDir), zap.String("savePath", savePath))
	return savePath
}

func prepareURLFileList(URL string) error {
	lst, err := getURL(URL)
	if err != nil {
		logger.Error("prepareURLFileList", zap.Error(err))
		return err
	}

	logger.Info("prepareURLFileList", zap.String("url content", lst))
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

func DownloadFile(URL string, localPath string, withProgress bool) error {
	timeStart := time.Now().Unix()
	logger.Info("DownloadFile:begin",
		zap.String("url", URL),
		zap.String("localPath", localPath),
		zap.String("time-start", strconv.FormatInt(timeStart, 10)))

	_, err := os.Stat(localPath)

	if err == nil {
		if IsOverwrite == true {
			err = os.Remove(localPath)
			if err != nil {
				logger.Error("DownloadFile: error(os-remove)", zap.String("cannot delete file", localPath), zap.Error(err))
			}
		} else {
			logger.Info("DownloadFile", zap.String("SKIP download", localPath), zap.String("is_overwirite", "true"))
			return nil
		}
	}

	resp, err := http.Get(URL)

	if err != nil {
		logger.Error("DownloadFile:error(http-get)", zap.Error(err))
		return err
	}
	defer resp.Body.Close()

	localPathTempName := strings.Join([]string{localPath, "downloading"}, ".")
	fileTemp, err := os.Create(localPathTempName)
	if err != nil {
		logger.Error("DownloadFile:error(os-create)", zap.String("cannot create file", localPathTempName), zap.Error(err))
		return err
	}

	defer fileTemp.Close()

	if withProgress == true {
		bar := config.SetBar(resp.ContentLength, "downloading").Bar
		_, err = io.Copy(io.MultiWriter(fileTemp, bar), resp.Body)
		bar.Finish()
	} else {
		_, err = io.Copy(fileTemp, resp.Body)
	}

	if err != nil {
		logger.Error("DownloadFile:error(io-copy)", zap.Error(err))
		return err
	}

	fileTemp.Close()

	err = os.Rename(localPathTempName, localPath)
	if err != nil {
		logger.Error("DownloadFile:error(os-rename)", zap.Error(err))
		return err
	}

	timeStop := time.Now().Unix()

	logger.Info("DownloadFile:ok",
		zap.String("proto", resp.Proto),
		zap.Int64("content-length", resp.ContentLength),
		zap.String("duration", strconv.FormatInt(timeStop-timeStart, 10)))

	return nil
}

func getURL(URL string) (cnt string, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		logger.Error("getURL", zap.String("url", URL), zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("getURL, reading body", zap.String("url", URL), zap.Error(err))
		return "", err
	}

	return string(body), nil
}

func makeDirs(p string) error {
	if _, err := os.Stat(p); err != nil {
		err := os.MkdirAll(p, os.ModePerm)
		if err != nil {
			logger.Error("cannot make dir:"+p, zap.Error(err))
			return err
		} else {
			os.Chmod(p, os.ModePerm)
		}
	}
	return nil
}

func Filepathify(fp string) string {
	var replacement string = "_"

	reControlCharsRegex := regexp.MustCompile("[\u0000-\u001f\u0080-\u009f]")

	reRelativePathRegex := regexp.MustCompile(`^\.+`)

	filenameReservedRegex := regexp.MustCompile(`[<>:"\\|?*\x00-\x1F]`)
	filenameReservedWindowsNamesRegex := regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])$`)

	// reserved word
	fp = filenameReservedRegex.ReplaceAllString(fp, replacement)

	// continue
	fp = reControlCharsRegex.ReplaceAllString(fp, replacement)
	fp = reRelativePathRegex.ReplaceAllString(fp, replacement)
	fp = filenameReservedWindowsNamesRegex.ReplaceAllString(fp, replacement)
	return fp
}

func StartDownload() error {
	if len(URLFileList) == 0 {
		logger.Info("SKIP download", zap.String("StartDownload", "URLFileList is empty"))
		return nil
	}

	for u, f := range URLFileList {
		if u == "" || f == "" {
			continue
		}

		DownloadFile(u, f, true)
	}

	return nil
}

func StartMultiDownload() error {
	if len(URLFileList) == 0 {
		return nil
	}

	wg := sync.WaitGroup{}
	var runningQueue int = 0

	for u, f := range URLFileList {
		logger.Info("start download", zap.String("url", u), zap.String("localPath", f))
		wg.Add(1)
		go func(u, f string) {
			DownloadFile(u, f, false)
			wg.Done()
		}(u, f)

		runningQueue++
		if runningQueue >= BatchSize {
			wg.Wait()
			runningQueue = 0
		}

	}
	wg.Wait()

	return nil
}
