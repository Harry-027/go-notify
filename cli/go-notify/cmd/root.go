package cmd

import (
	"github.com/spf13/cobra"
)

// Cobra root command ...
var RootCmd = &cobra.Command{
	Use: "go-notify",
	Short: "A CLI for an email marketing tool: go-notify that lets users to register & schedule/send custom mails to their clients(leads/prospects). \n" +
		"*) Sign-up \n" +
		"*) Login \n" +
		"*) Update password \n" +
		"*) Logout \n" +
		"*) Forgot password" +
		"*) Subscribe \n" +
		"*) Add Clients \n" +
		"*) Get Clients \n" +
		"*) Add Template \n" +
		"*) Get Template \n" +
		"*) Get Subscription details \n" +
		"*) Get Users \n" +
		"*) Add client variable for a template \n" +
		"*) Get client's variable \n" +
		"*) Delete a client \n" +
		"*) Delete a template \n" +
		"*) Get subscription details \n" +
		"*) Send Mail \n" +
		"*) Schedule Mail \n" +
		"*) Cancel scheduled mail \n" +
		"*) Delete Account",
}