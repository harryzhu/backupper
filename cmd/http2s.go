/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"sqlconf"

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
		allowiplist := sqlconf.StringToSlice(config.ToString("http2s_allow_ip_list"))
		blockiplist := sqlconf.StringToSlice(config.ToString("http2s_block_ip_list"))
		defaultallow := config.ToBool("http2s_default_allow")

		h2s := config.SetH2Server().H2Server
		h2s.WithStaticRootDir(config.ToString("http2s_static_root_dir"))
		for _, aip := range allowiplist {
			h2s.WithAllowIP(aip)
		}
		for _, bip := range blockiplist {
			h2s.WithBlockIP(bip)
		}
		h2s.WithDefaultAllow(defaultallow)
		h2s.WithIP("127.0.0.1")
		h2s.WithPort(config.ToInt("http2s_port"))
		h2s.WithTLS(config.ToString("http2s_tls_cert"), config.ToString("http2s_tls_key"))
		h2s.StartServer()
	},
}

func init() {
	rootCmd.AddCommand(http2sCmd)

}
