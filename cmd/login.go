/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/quexten/goldwarden/ipc"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Starts the login process for Bitwarden",
	Long: `Starts the login process for Bitwarden. 
	You will be prompted to enter your password, and confirm your second factor if you have one.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := ipc.DoLoginRequest{}
		email, _ := cmd.Flags().GetString("email")
		request.Email = email
		passwordless, _ := cmd.Flags().GetBool("passwordless")
		request.Passwordless = passwordless

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Logged in")
			} else {
				println("Login failed: " + result.(ipc.ActionResponse).Message)
			}
		default:
			println("Wrong IPC response type for login")
		}
	},
}

func init() {
	vaultCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("email", "", "")
	loginCmd.MarkFlagRequired("email")
	loginCmd.PersistentFlags().Bool("passwordless", false, "")
}
