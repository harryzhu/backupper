/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Force      bool
	initConfig map[string]string
)

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
		fmt.Println("reset called")
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
