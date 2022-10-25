/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// http2Cmd represents the http2 command
var http2sCmd = &cobra.Command{
	Use:   "http2s",
	Short: "A brief description of your command",
	Long:  `-`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http2 server called")
		//var StaticRootDir string = "/Volumes/HDD4/downloads/"

		h2s := config.SetH2Server().H2Server
		h2s.WithStaticRootDir(config.ToString("http2s_static_root_dir"))
		h2s.WithIP("127.0.0.1")
		h2s.WithPort(config.ToInt("http2s_port"))
		h2s.WithTLS(config.ToString("http2s_tls_cert"), config.ToString("http2s_tls_key"))
		h2s.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(http2sCmd)

}
