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

// AddClient command ...
var addClients = &cobra.Command{
	Use:   "addClient",
	Short: "Add client details (Leads/Prospects)",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "POST"
		urlPath := "/api/clients"
		callType := user
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

// GetClients command ...
var getClients = &cobra.Command{
	Use:   "getClient",
	Short: "Get the registered client's detail (Leads/Prospects)",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "GET"
		urlPath := "/api/clients"
		callType := user

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		if statusCode == "200 OK" {
			headers := []string{"ClientId", "Email", "Name", "Phone", "Preference"}
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

// DelClient command ...
var delClient = &cobra.Command{
	Use:   "delClient",
	Short: "Remove the client from system ..",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var mailId string
		methodType := "DELETE"
		urlPath := "/api/deleteClient"
		callType := user

		fmt.Print("Enter client's mail-id: ")
		fmt.Scanln(&mailId)

		urlPath = fmt.Sprintf("%s/%s", urlPath, mailId)
		config := readConfig()
		_, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		err := json.Unmarshal(body, &apiResponse)
		if err == nil {
			fmt.Println("Status: ", apiResponse.Status)
			fmt.Println("Message: ", apiResponse.Message)
		}
	},
}

// Cobra Clients command ...
func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	payloadPath := filepath.Join(home, defaultPayloadPath)
	addClients.Flags().String(PayloadPath, payloadPath, payloadPathDesc)
	RootCmd.AddCommand(addClients)
	RootCmd.AddCommand(getClients)
	RootCmd.AddCommand(delClient)
}
