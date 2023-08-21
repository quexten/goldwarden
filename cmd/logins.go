/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc"
	"github.com/spf13/cobra"
)

var baseLoginCmd = &cobra.Command{
	Use:   "logins",
	Short: "Commands for managing logins.",
	Long:  `Commands for managing logins.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var getLoginCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a login in your vault",
	Long:  `Gets a login in your vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		loginIfRequired()

		uuid, _ := cmd.Flags().GetString("uuid")
		name, _ := cmd.Flags().GetString("name")
		username, _ := cmd.Flags().GetString("username")
		fullOutput, _ := cmd.Flags().GetBool("full")

		resp, err := commandClient.SendToAgent(ipc.GetLoginRequest{
			Name:     name,
			Username: username,
			UUID:     uuid,
		})
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
			return
		}

		switch resp.(type) {
		case ipc.GetLoginResponse:
			response := resp.(ipc.GetLoginResponse)
			if fullOutput {
				fmt.Println(response.Result)
			} else {
				fmt.Println(response.Result.Password)
			}
			break
		case ipc.ActionResponse:
			println("Error: " + resp.(ipc.ActionResponse).Message)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(baseLoginCmd)
	baseLoginCmd.AddCommand(getLoginCmd)
	getLoginCmd.PersistentFlags().String("name", "", "")
	getLoginCmd.PersistentFlags().String("username", "", "")
	getLoginCmd.PersistentFlags().String("uuid", "", "")
	getLoginCmd.PersistentFlags().Bool("full", false, "")
}
