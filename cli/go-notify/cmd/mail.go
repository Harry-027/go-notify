package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// SendMail command ...
var sendMail = &cobra.Command{
	Use:   "sendMail",
	Short: "send mail to the clients",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		callType := user
		methodType := "POST"
		urlPath := "/api/sendMail"
		payloadPath, err := cmd.Flags().GetString(PayloadPath)
		data, err := ioutil.ReadFile(payloadPath)
		if err != nil {
			log.Println("File reading error", err)
			return
		}

		config := readConfig()
		_, body := makeCallToServer(methodType, callType, urlPath, config.Token, data)

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("An error occurred while data unmarshal", err)
		}

		fmt.Println("Status: ", apiResponse.Status)
		fmt.Println("Message: ", apiResponse.Message)
	},
}

// ScheduleMail command ...
var scheduleMail = &cobra.Command{
	Use:   "scheduleMail",
	Short: "schedule mail for the clients as per their preference (daily, weekly, monthly)",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		callType := user
		methodType := "POST"
		urlPath := "/api/scheduleMail"
		payloadPath, err := cmd.Flags().GetString(PayloadPath)
		data, err := ioutil.ReadFile(payloadPath)
		if err != nil {
			log.Println("File reading error", err)
			return
		}

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, data)

		if statusCode == "200 OK" {
			headers := []string{"ClientId", "JobId"}
			stringSlice := dataConversion(body, headers)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(headers)
			for _, v := range stringSlice {
				table.Append(v)
			}
			table.SetAlignment(1)
			table.Render()
		} else {
			err = json.Unmarshal(body, &apiResponse)
			if err != nil {
				fmt.Println("An error occurred while data unmarshal", err)
			}
			fmt.Println("Status: ", apiResponse.Status)
			fmt.Println("Message: ", apiResponse.Message)
		}
	},
}

// cancelScheduledMail command
var cancelScheduledMail = &cobra.Command{
	Use:   "cancelScheduleMail",
	Short: "Cancel the scheduled mails",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var jobId int
		methodType := "POST"
		urlPath := "/api/deleteScheduleMail"
		callType := user

		fmt.Print("Enter job id: ")
		fmt.Scanln(&jobId)

		data, err := json.Marshal(map[string]int{
			"jobId": jobId,
		})

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, data)

		err = json.Unmarshal(body, &apiResponse)
		if err == nil {
			fmt.Println("StatusCode: ", statusCode)
			fmt.Println("Status: ", apiResponse.Status)
			fmt.Println("Message: ", apiResponse.Message)
		} else {
			fmt.Println("An error occurred during response unmarshal", err.Error())
		}
	},
}

// GetTemplates command ...
var getAuditLog = &cobra.Command{
	Use:   "getAuditLog",
	Short: "Fetch the audit log for sent mails",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "GET"
		urlPath := "/api/checkAuditLog"
		callType := user

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		if statusCode == "200 OK" {
			headers := []string{"SentOn", "To", "template"}
			stringSlice := dataConversion(body, headers)
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader(headers)
			for _, v := range stringSlice {
				table.Append(v)
			}
			table.SetAlignment(1)
			table.Render()
		} else {
			err := json.Unmarshal(body, &apiResponse)
			if err == nil {
				fmt.Println("Status: ", apiResponse.Status)
				fmt.Println("Message: ", apiResponse.Message)
			}
		}
	},
}

func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	payloadPath := filepath.Join(home, defaultPayloadPath)
	sendMail.Flags().String(PayloadPath, payloadPath, payloadPathDesc)
	scheduleMail.Flags().String(PayloadPath, payloadPath, payloadPathDesc)
	RootCmd.AddCommand(sendMail)
	RootCmd.AddCommand(scheduleMail)
	RootCmd.AddCommand(cancelScheduledMail)
	RootCmd.AddCommand(getAuditLog)
}
