package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Commands for managing sends",
	Long:  `Commands for managing sends.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var sendCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Uploads a Bitwarden send.",
	Long:  `Uploads a Bitwarden send.`,
	Run: func(cmd *cobra.Command, args []string) {
		loginIfRequired()
		name, _ := cmd.Flags().GetString("name")
		text, _ := cmd.Flags().GetString("text")

		result, err := commandClient.SendToAgent(messages.CreateSendRequest{
			Name: name,
			Text: text,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.CreateSendResponse:
			fmt.Println("Send created: " + result.(messages.CreateSendResponse).URL)
			break
		case messages.ActionResponse:
			println("Error: " + result.(messages.ActionResponse).Message)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.AddCommand(sendCreateCmd)
	sendCreateCmd.Flags().StringP("name", "n", "", "Name of the send")
	sendCreateCmd.Flags().StringP("text", "t", "", "Text of the send")
}
