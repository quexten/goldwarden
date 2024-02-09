/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Starts the login process for Bitwarden",
	Long: `Starts the login process for Bitwarden. 
	You will be prompted to enter your password, and confirm your second factor if you have one.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.DoLoginRequest{}
		email, _ := cmd.Flags().GetString("email")
		if email == "" {
			fmt.Println("Error: No email specified")
			return
		}

		request.Email = email
		passwordless, _ := cmd.Flags().GetBool("passwordless")
		request.Passwordless = passwordless

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				fmt.Println("Logged in")
			} else {
				fmt.Println("Login failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type for login")
		}
	},
}

func init() {
	vaultCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("email", "", "")
	loginCmd.MarkFlagRequired("email")
	loginCmd.PersistentFlags().Bool("passwordless", false, "")
}
