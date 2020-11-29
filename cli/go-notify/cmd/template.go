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

// AddTemplate command ...
var addTemplate = &cobra.Command{
	Use: "addTemplate",
	Short: "Add custom HTML template to be used for sending mails to clients.",
	Long: "Add custom HTML template to be used for sending mails to clients. You can also introduce template variables in mails." +
		  "Syntax to add the template variable in template is: {{ variableName }}",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "POST"
		urlPath := "/api/addTemplate"
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

// GetTemplates command ...
var getTemplates = &cobra.Command{
	Use:   "getTemplates",
	Short: "Fetch the registered custom HTML templates",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "GET"
		urlPath := "/api/getTemplates"
		callType := user

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		if statusCode == "200 OK" {
			headers := []string{"name", "subject", "body"}
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

// UpdateTemplate command ...
var updateTemplate = &cobra.Command{
	Use:   "updateTemplate",
	Short: "Update the mail template",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var tempName string
		methodType := "PUT"
		urlPath := "/api/updateTemplate"
		callType := user

		payloadPath, err := cmd.Flags().GetString(PayloadPath)
		data, err := ioutil.ReadFile(payloadPath)
		if err != nil {
			log.Println("File reading error", err)
			return
		}

		fmt.Print("Enter template name: ")
		fmt.Scanln(&tempName)
		urlPath = fmt.Sprintf("%s/%s", urlPath, tempName)
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

// DeleteTemplate command ...
var delTemplate = &cobra.Command{
	Use:   "delTemplate",
	Short: "Remove the registered custom HTML template from system ..",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var name string
		methodType := "DELETE"
		urlPath := "/api/deleteTemplate"
		callType := user

		fmt.Print("Enter template name: ")
		fmt.Scanln(&name)

		urlPath = fmt.Sprintf("%s/%s", urlPath, name)
		config := readConfig()
		_, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		err := json.Unmarshal(body, &apiResponse)
		if err == nil {
			fmt.Println("Status: ", apiResponse.Status)
			fmt.Println("Message: ", apiResponse.Message)
		}
	},
}

// AddTemplate command ...
var addTemplateVar = &cobra.Command{
	Use:   "addTemplateVar",
	Short: "Add client specific details for a template variable",
	Long: "The custom HTML template may contain variables in the format - {{ variableName }}. Those variables" +
		"can be replaced by client specific detail. This command helps to register the client specific detail for a given variable",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "POST"
		urlPath := "/api/clientDetails"
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

// GetClientVar command ...
var getClientVar = &cobra.Command{
	Use:   "getTemplateVar",
	Short: "Fetch the values for template variables for a specific client.",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var clientId string
		methodType := "GET"
		urlPath := "/api/clientDetails"
		callType := user
		fmt.Print("Enter Client Id: ")
		fmt.Scanln(&clientId)
		urlPath = fmt.Sprintf("%s/%s", urlPath, clientId)
		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		if statusCode == "200 OK" {
			headers := []string{"Key", "Value"}
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

// DelClientVar command ...
var delClientVar = &cobra.Command{
	Use:   "delClientVar",
	Short: "Delete template variable value for a given client.",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var mailId string
		var key string
		methodType := "DELETE"
		urlPath := "/api/clientDetails"
		callType := user
		fmt.Print("Enter Client MailId: ")
		fmt.Scanln(&mailId)
		fmt.Print("Enter template variable name: ")
		fmt.Scanln(&key)

		data, err := json.Marshal(map[string]string{
			"key":          key,
			"clientMailID": mailId,
		})
		if err != nil {
			fmt.Println("An error occurred during marshal: ", err)
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

// Cobra Clients command ...
func init() {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	templatePayloadPath := filepath.Join(home, defaultPayloadPath)
	addTemplate.Flags().String(PayloadPath, templatePayloadPath, payloadPathDesc)
	updateTemplate.Flags().String(PayloadPath, templatePayloadPath, payloadPathDesc)
	addTemplateVar.Flags().String(PayloadPath, templatePayloadPath, payloadPathDesc)
	RootCmd.AddCommand(addTemplate)
	RootCmd.AddCommand(getTemplates)
	RootCmd.AddCommand(delTemplate)
	RootCmd.AddCommand(updateTemplate)
	RootCmd.AddCommand(addTemplateVar)
	RootCmd.AddCommand(getClientVar)
	RootCmd.AddCommand(delClientVar)
}
