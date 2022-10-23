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
		h2s.WithStaticRootDir("/Volumes/HDD4/downloads/").WithAddress(":8080")
		h2s.WithTLS("../../../cert/cert.pem", "../../../cert/priv.key")
		h2s.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(http2sCmd)

}
