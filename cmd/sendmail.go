/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	File    string
	Subject string
)

// sendmailCmd represents the sendmail command
var sendmailCmd = &cobra.Command{
	Use:   "sendmail",
	Short: "A brief description of your command",
	Long:  `-`,
	Run: func(cmd *cobra.Command, args []string) {

		if _, err := os.Stat(File); err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadFile(File)
		if err != nil {
			log.Fatal(err)
		}

		config.Mail.WithMessage(Subject, string(content)).SendMailStartTLS()
	},
}

func init() {
	rootCmd.AddCommand(sendmailCmd)
	sendmailCmd.Flags().StringVar(&File, "file", "", "local file of the mail body")
	sendmailCmd.Flags().StringVar(&Subject, "subject", "", "mail subject")

	sendmailCmd.MarkFlagRequired("file")
	sendmailCmd.MarkFlagRequired("subject")
}
