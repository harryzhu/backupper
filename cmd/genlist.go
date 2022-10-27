/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	//"fmt"
	"sqlconf"

	"github.com/spf13/cobra"
)

// genlistCmd represents the genlist command
var genlistCmd = &cobra.Command{
	Use:   "genlist",
	Short: "A brief description of your command",
	Long:  `-`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("genlist called")
		sqlconf.GenFileListByDir(config.ToString("genlist_root_dir"),
			config.ToString("genlist_url_prefix"),
			config.ToString("genlist_out_file"))

	},
}

func init() {
	rootCmd.AddCommand(genlistCmd)

}
