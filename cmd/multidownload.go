/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var multidownloadCmd = &cobra.Command{
	Use:   "multidownload",
	Short: "A brief description of your command",
	Long:  `-`,
	PreRun: func(cmd *cobra.Command, args []string) {
		config.RequiredKeys([]string{"url_list", "dir_save_root"})
	},
	Run: func(cmd *cobra.Command, args []string) {
		StartMultiDownload()
	},
}

func init() {
	rootCmd.AddCommand(multidownloadCmd)
	multidownloadCmd.Flags().IntVar(&BatchSize, "batch", 3, "--batch=3|5|10")
}
