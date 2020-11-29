package cmd

import (
	"github.com/spf13/cobra"
)

// Cobra root command ...
var RootCmd = &cobra.Command{
	Use: "go-notify",
	Short: "A CLI for an email marketing tool: go-notify that lets users to register & schedule/send custom mails to their clients(leads/prospects). \n" +
		"Few command's input would be a json file. Please refer payloadSample.js or swagger api document to determine the input payload for those commands",
}
