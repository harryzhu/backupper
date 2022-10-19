/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.RequiredKeys([]string{"url_list", "dir_save_root"})
	},
	Run: func(cmd *cobra.Command, args []string) {
		StartDownload()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolVar(&IsOverwrite, "overwrite", false, "--overwrite=false|true")
}
