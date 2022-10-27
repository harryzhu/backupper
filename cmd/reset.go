/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"sqlconf"
	"strings"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	IsForce     bool
	IsLocal     bool
	resetFile   string
	resetConfig map[string]string
	config      *sqlconf.Conf
	logger      *zap.Logger
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
		resetConfig = make(map[string]string, 10)
		// default config
		resetFile = strings.Replace(config.DBFile, "conf.db", "reset.txt", 1)
		rf, err := ioutil.ReadFile(resetFile)
		if err != nil {
			logger.Error("cannot read resetfile", zap.Error(err))
		}
		strrf := strings.ReplaceAll(string(rf), "\r\n", "\n")
		lines := strings.Split(strrf, "\n")
		var arrLine []string
		key := ""
		val := ""
		for _, line := range lines {
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, " ")
			if line == "" || strings.Index(line, "#") == 0 || strings.Index(line, "=") <= 0 {
				continue
			}
			arrLine = strings.Split(line, "=")
			if len(arrLine) != 2 {
				continue
			}
			key = strings.Trim(arrLine[0], " ")
			val = strings.Trim(arrLine[1], " ")
			if arrLine[0] == "" {
				continue
			}
			resetConfig[key] = val

			//fmt.Println(key, "=", val)
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset, using", resetFile)
		if IsForce == true {
			config.LoadData(resetConfig)
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
