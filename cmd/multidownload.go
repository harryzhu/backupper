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
	Run: func(cmd *cobra.Command, args []string) {
		StartMultiDownload()
	},
}

func init() {
	rootCmd.AddCommand(multidownloadCmd)

}
