package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc/messages"
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

		resp, err := commandClient.SendToAgent(messages.GetLoginRequest{
			Name:     name,
			Username: username,
			UUID:     uuid,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch resp.(type) {
		case messages.GetLoginResponse:
			response := resp.(messages.GetLoginResponse)
			if fullOutput {
				fmt.Println(response.Result)
			} else {
				fmt.Println(response.Result.Password)
			}
			break
		case messages.ActionResponse:
			println("Error: " + resp.(messages.ActionResponse).Message)
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
