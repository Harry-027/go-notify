package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
)

// SubsDetail command ...
var subsDetail = &cobra.Command{
	Use:   "subsDetail",
	Short: "Get the current quota left as per the subscription plan",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		callType := user
		methodType := "GET"
		urlPath := "/api/subscriptionDetails"
		config := readConfig()
		_, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)

		err := json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("An error occurred while data unmarshal", err)
		}

		fmt.Println("Status: ", apiResponse.Status)
		fmt.Println("Message: ", apiResponse.Message)
	},
}

// subscribe command ...
var subs = &cobra.Command{
	Use:   "subscribe",
	Short: "Subscribe to plan - (gold, silver, platinum)",
	Run: func(cmd *cobra.Command, args []string) {
		var plan, paymentType string
		var apiResponse ApiResponse
		callType := user
		methodType := "POST"
		urlPath := "/api/subscribe"
		fmt.Print("Select your plan (gold, silver, platinum): ")
		fmt.Scanln(&plan)
		fmt.Print("Enter the payment mode: (credit, debit, netBanking): ")
		fmt.Scanln(&paymentType)
		data, err := json.Marshal(map[string]string{
			"subscriptionType": plan,
			"paymentType":      paymentType,
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

// Cobra subscription command ...
func init() {
	RootCmd.AddCommand(subs)
	RootCmd.AddCommand(subsDetail)
}
