/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"errors"
	"fmt"
	//"io/ioutil"
	"os"
	//"path/filepath"
	"sqlconf"
	"strconv"

	//"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ts_now int64         = time.Now().Unix()
	config *sqlconf.Conf = sqlconf.Config

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

		prepareURLFileList(config.ToString("url_list"))

		fmt.Println("======")
		//fmt.Println(URLFileList)
	},
	Run: func(cmd *cobra.Command, args []string) {

	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.Close()

		//config.AlwaysPostRun()

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

	logger = config.SetLogger().Logger.ZapLogger
	logger.Info("Thank you for choosing " + config.ToString("app_name"))

	config.SetMail()

	DirSaveRoot = config.ToString("dir_save_root")
}
