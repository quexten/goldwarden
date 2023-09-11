package cmd

import (
	"github.com/quexten/goldwarden/ipc"
	"github.com/spf13/cobra"
)

var setApiUrlCmd = &cobra.Command{
	Use:   "set-api-url",
	Short: "Set the api url",
	Long:  `Set the api url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		url := args[0]
		request := ipc.SetApiURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting api url failed: " + result.(ipc.ActionResponse).Message)
			}
		default:
			println("Wrong IPC response type")
		}

	},
}

var setIdentityURLCmd = &cobra.Command{
	Use:   "set-identity-url",
	Short: "Set the identity url",
	Long:  `Set the identity url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		url := args[0]
		request := ipc.SetIdentityURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting identity url failed: " + result.(ipc.ActionResponse).Message)
			}
		default:
			println("Wrong IPC response type")
		}

	},
}

var setNotificationsURLCmd = &cobra.Command{
	Use:   "set-notifications-url",
	Short: "Set the notifications url",
	Long:  `Set the notifications url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		url := args[0]
		request := ipc.SetNotificationsURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting notifications url failed: " + result.(ipc.ActionResponse).Message)
			}
		default:
			println("Wrong IPC response type")
		}

	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration",
	Long:  `Manage the configuration.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setApiUrlCmd)
	configCmd.AddCommand(setIdentityURLCmd)
	configCmd.AddCommand(setNotificationsURLCmd)
}
