/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sqlconf"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ts_now          int64         = time.Now().Unix()
	config          *sqlconf.Conf = new(sqlconf.Conf)
	initConfig      map[string]string
	globalTimeStart int64
	globalTimeStop  int64
	logger          *zap.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backupper",
	Short: "A brief description of your application",
	Long:  `-`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		globalTimeStart = time.Now().Unix()
		config.Refresh().Print()

		logger.Info("===== start =====", zap.String("time", strconv.FormatInt(globalTimeStart, 10)))
		//logger.Error("ddd", zap.Error(errors.New("error test")))

		prepareURLFileList(config.ToString("url_list"))

		fmt.Println("======")
		//fmt.Println(URLFileList)
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.Close()
		curErrorLog := sqlconf.CurrentErrorLog
		logger.Info("curErrorLog", zap.String("curErrorLog", curErrorLog))

		if fi, err := os.Stat(curErrorLog); err == nil {
			if fi.Size() > 16 {
				cnt, err := ioutil.ReadFile(curErrorLog)
				if err == nil {
					mSubject := "[ERROR].[RUNNING]:" + filepath.Base(curErrorLog)
					mBody := strings.Join([]string{curErrorLog, "<br/><br/>", "<pre>", string(cnt), "</pre>"}, "")
					config.Mail.WithMessage(mSubject, mBody).SendMailStartTLS()
				}

			}
		}

		globalTimeStop = time.Now().Unix()
		logger.Info("===== end(total duration) =====", zap.String("second", strconv.FormatInt(globalTimeStop-globalTimeStart, 10)))

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	initConfig = make(map[string]string, 10)

	initConfig["app_first_run"] = strconv.FormatInt(ts_now, 10)
	initConfig["app_conf_update"] = strconv.FormatInt(ts_now, 10)
	initConfig["app_name"] = "backupper"
	initConfig["app_author"] = "harryzhu"
	initConfig["app_license"] = "MIT"
	initConfig["app_version"] = "1.0.0"
	initConfig["app_data_dir"] = "./data"
	initConfig["app_logs_dir"] = "./logs"
	initConfig["app_temp_dir"] = "./temp"

	// err
	initConfig["err_empty"] = "value cannot be empty"

	initConfig["url_list"] = ""
	initConfig["dir_save_root"] = ""

	bootConfig()
}

func bootConfig() {
	cFile := "./conf.db"
	firstRun := false

	_, err := os.Stat(cFile)
	if err != nil {
		firstRun = true
	}
	config.Open(cFile)

	if firstRun == true {
		config.LoadData(initConfig)
	}

	config.Refresh()
	config.FatalEmpty([]string{"app_name", "app_logs_dir"})
	config.FatalEmpty([]string{"url_list", "dir_save_root"})

	config.SetLogger(config.ToString("app_logs_dir"), config.ToString("app_name"))
	config.SetMail()

	logger = config.Logger
	logger.Info("Thank you for choosing " + config.ToString("app_name"))

	DirSaveRoot = config.ToString("dir_save_root")

}
