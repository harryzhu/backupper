/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sqlconf"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	Force      bool
	initConfig map[string]string
	config     *sqlconf.Conf = sqlconf.Config
	logger     *zap.Logger
)

func bootConfig() {
	logger = config.SetLogger().Logger.ZapLogger
	config.SetMail()

}

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "A brief description of your command",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig = make(map[string]string, 10)
		// default config
		initConfig["url_list"] = "http://localhost/download_v2.txt"
		initConfig["dir_save_root"] = "/Volumes/SSD2/temp/"

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset")
		if Force == true {
			config.LoadData(initConfig)
			config.Refresh().Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolVar(&Force, "force", false, "--force=true")
	resetCmd.MarkFlagRequired("force")
}
