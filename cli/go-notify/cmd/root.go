package cmd

import (
	"github.com/spf13/cobra"
)

// Cobra root command ...
var RootCmd = &cobra.Command{
	Use: "go-notify",
	Short: "A CLI for an email automation solution: go-notify that facilitate users to register, send & schedule custom HTML mails for their clients. \n" +
		"Few command's input would be a json file. Please refer payloadSample.js or swagger api document to determine the input payload for those commands",
}
