package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

// Signin command ...
var signIn = &cobra.Command{
	Use:   "signin",
	Short: "Login as a User",
	Run: func(cmd *cobra.Command, args []string) {
		var mailId, pwd string
		var loginResponse LoginResponse
		callType := auth
		methodType := "POST"
		urlPath := "/auth/login"
		fmt.Print("Enter your mail-id: ")
		fmt.Scanln(&mailId)
		fmt.Print("Enter password: ")
		fmt.Scanln(&pwd)

		data, err := json.Marshal(map[string]string{
			"email":    mailId,
			"password": pwd,
		})
		if err != nil {
			fmt.Println("An error occurred during marshal: ", err)
		}

		_, body := makeCallToServer(methodType, callType, urlPath, "", data)

		err = json.Unmarshal(body, &loginResponse)
		if err != nil {
			fmt.Println("An error occurred while data unmarshal", err)
		}

		if loginResponse.Status == "Success" {
			config := Config{Token: loginResponse.Token}
			err := saveConfig(config)
			if err != nil {
				fmt.Println("Internal server error occurred ..")
			} else {
				fmt.Println("Message: ", loginResponse.Message)
			}
		} else {
			fmt.Println("Message: Login failed !!")
		}
	},
}

// Logout command ...
var logout = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the session",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		callType := user
		methodType := "POST"
		urlPath := "/privacy/logout"

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

// Forgot password command ...
var fgt = &cobra.Command{
	Use:   "resetPwd",
	Short: "Forgot the password ? Reset now !!",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		var mail string
		callType := auth
		methodType := "POST"
		urlPath := "/auth/forgotPassword"

		fmt.Print("Enter your registered mailId: ")
		fmt.Scanln(&mail)

		data, err := json.Marshal(map[string]string{
			"email": mail,
		})
		if err != nil {
			fmt.Println("An error occurred during marshal: ", err)
		}

		_, body := makeCallToServer(methodType, callType, urlPath, "", data)
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("An error occurred while data unmarshal", err)
		}
		fmt.Println("Message: ", apiResponse.Message)
	},
}

// GetUsers command ...
var getUsers = &cobra.Command{
	Use:   "getUsers",
	Short: "Get all the registered users in the system. Your role should be an admin to get the desired results",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "GET"
		urlPath := "/api/users"
		callType := user

		config := readConfig()
		statusCode, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		if statusCode == "200 OK" {
			headers := []string{"Email", "NotificationCounter", "Role", "Subscription"}
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

// Signup command ...
var signCmd = &cobra.Command{
	Use:   "signup",
	Short: "Signup to register as a User. You may decide your role here as - Admin or User",
	Run: func(cmd *cobra.Command, args []string) {
		var mailId, isAdmin, role, pwd, cnfPwd string
		var apiResponse ApiResponse
		callType := auth
		methodType := "POST"
		urlPath := "/auth/signup"

		fmt.Print("Enter your mail-id: ")
		fmt.Scanln(&mailId)
		fmt.Print("Are you an admin(Y/N): ")
		fmt.Scanln(&isAdmin)
		fmt.Print("Enter password: ")
		fmt.Scanln(&pwd)
		fmt.Print("Confirm password: ")
		fmt.Scanln(&cnfPwd)

		if isAdmin == "Y" || isAdmin == "y" {
			role = "admin"
		} else {
			role = "user"
		}

		data, err := json.Marshal(map[string]string{
			"email":            mailId,
			"password":         pwd,
			"confirm_password": cnfPwd,
			"role":             role,
		})
		if err != nil {
			fmt.Println("An error occurred during marshal: ", err)
		}

		_, body := makeCallToServer(methodType, callType, urlPath, "", data)

		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println("An error occurred while data unmarshal", err)
		}

		fmt.Println("Status: ", apiResponse.Status)
		fmt.Println("Message: ", apiResponse.Message)
	},
}

// UpdatePassword command ...
var updatePassword = &cobra.Command{
	Use:   "updatePwd",
	Short: "Update your password !!",
	Run: func(cmd *cobra.Command, args []string) {
		var pwd string
		var apiResponse ApiResponse
		callType := user
		methodType := "POST"
		urlPath := "/privacy/updatePassword"
		fmt.Print("Enter new password: ")
		fmt.Scanln(&pwd)

		data, err := json.Marshal(map[string]string{
			"password": pwd,
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

// DeleteAccount command ...
var deleteAccount = &cobra.Command{
	Use:   "delAcc",
	Short: "Remove your account permanently. Cautious : All data would be lost ..(Irreversible Action)",
	Run: func(cmd *cobra.Command, args []string) {
		var apiResponse ApiResponse
		methodType := "DELETE"
		urlPath := "/privacy/deleteAccount"
		callType := user

		config := readConfig()
		_, body := makeCallToServer(methodType, callType, urlPath, config.Token, nil)
		err := json.Unmarshal(body, &apiResponse)
		if err == nil {
			fmt.Println("Status: ", apiResponse.Status)
			fmt.Println("Message: ", apiResponse.Message)
		}
	},
}

// Cobra auth commands ...
func init() {
	RootCmd.AddCommand(signCmd)
	RootCmd.AddCommand(signIn)
	RootCmd.AddCommand(logout)
	RootCmd.AddCommand(fgt)
	RootCmd.AddCommand(updatePassword)
	RootCmd.AddCommand(getUsers)
	RootCmd.AddCommand(deleteAccount)
}
