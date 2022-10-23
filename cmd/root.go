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
	//"strconv"

	//"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	globalTimeStart int64
	globalTimeStop  int64
	AutoRunAlways   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "backupper",
	Short: "A brief description of your application",
	Long:  `-`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		globalTimeStart = time.Now().Unix()
		config.Refresh().Print()

		logger.Info("===== start =====", zap.Int64("time", globalTimeStart))

		fmt.Println("======")
		//fmt.Println(URLFileList)
	},
	Run: func(cmd *cobra.Command, args []string) {
		config.Refresh().Print()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		config.Close()

		if AutoRunAlways == true {
			config.AlwaysPostRun()
		}

		globalTimeStop = time.Now().Unix()
		logger.Info("===== end =====", zap.Int64("duration", globalTimeStop-globalTimeStart))
		fmt.Println("***** END *****")
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
	rootCmd.PersistentFlags().BoolVar(&AutoRunAlways, "auto-run-always", false, "--auto-run-always=true|false")

	fmt.Println("***** BEGIN *****")
	bootConfig()

	logger.Info("Thank you for choosing " + config.ToString("app_name"))
	DirSaveRoot = config.ToString("dir_save_root")
}
