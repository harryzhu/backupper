/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var (
	IsHttp2 bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.RequiredKeys([]string{"url_list", "dir_save_root"})
		prepareURLFileList(config.ToString("url_list"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		if BatchSize > 1 {
			runMultiDownload()
		} else {
			runDownload()
		}

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().IntVar(&BatchSize, "batch", 1, "--batch=3|5|10")
	downloadCmd.Flags().BoolVar(&IsOverwrite, "overwrite", false, "--overwrite=false|true")
	downloadCmd.Flags().BoolVar(&IsHttp2, "http2", false, "--http2=false|true")
}
