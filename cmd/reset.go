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
	IsForce    bool
	IsLocal    bool
	initConfig map[string]string
	config     *sqlconf.Conf
	logger     *zap.Logger
)

func bootConfig() {
	config = sqlconf.NewConf("backupper", "./logs")
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
		//download
		if IsLocal {
			initConfig["download_url_list"] = "https://dla.harryzhu.plus:8080/download_v2.txt"
			initConfig["download_save_dir"] = "d:/backup/"
		} else {
			initConfig["download_url_list"] = "https://dla.harryzhu.plus/download_v2.txt"
			initConfig["download_save_dir"] = "./data"
		}

		//http2s
		if IsLocal {
			initConfig["http2s_ip"] = "127.0.0.1"
			initConfig["http2s_port"] = "8080"
			initConfig["http2s_static_root_dir"] = "d:/static/"
			initConfig["http2s_allow_ip_list"] = ""
			initConfig["http2s_block_ip_list"] = ""
			initConfig["http2s_default_allow"] = ""
			initConfig["http2s_tls_cert"] = "../../../cert/dla.harryzhu.plus.pem"
			initConfig["http2s_tls_key"] = "../../../cert/dla.harryzhu.plus.key"
			initConfig["http2s_enable_control"] = "1"
			initConfig["http2s_enable_reverse_proxy"] = "1"
			initConfig["http2s_reverse_proxy_url"] = ""
		} else {
			initConfig["http2s_ip"] = "0.0.0.0"
			initConfig["http2s_port"] = "443"
			initConfig["http2s_static_root_dir"] = "./"
			initConfig["http2s_allow_ip_list"] = ""
			initConfig["http2s_block_ip_list"] = ""
			initConfig["http2s_default_allow"] = ""
			initConfig["http2s_tls_cert"] = "../cert/dla.harryzhu.plus.pem"
			initConfig["http2s_tls_key"] = "../cert/dla.harryzhu.plus.key"
			initConfig["http2s_enable_control"] = "0"
			initConfig["http2s_enable_reverse_proxy"] = "0"
			initConfig["http2s_reverse_proxy_url"] = ""
		}

		// genlist
		if IsLocal {
			initConfig["genlist_root_dir"] = "d:/static/"
			initConfig["genlist_url_prefix"] = "https://dla.harryzhu.plus:8080"
			initConfig["genlist_out_file"] = "d:/static/download_v2.txt"
		} else {
			initConfig["genlist_root_dir"] = "./"
			initConfig["genlist_url_prefix"] = "https://dla.harryzhu.plus"
			initConfig["genlist_out_file"] = "./download_v2.txt"
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset")
		if IsForce == true {
			config.LoadData(initConfig)
			config.Refresh().Print()
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.Flags().BoolVar(&IsForce, "force", false, "--force=true")
	resetCmd.Flags().BoolVar(&IsLocal, "local", false, "--local=true")
	resetCmd.MarkFlagRequired("force")
}
