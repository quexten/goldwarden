package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc/messages"
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
		request := messages.SetApiURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting api url failed: " + result.(messages.ActionResponse).Message)
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
		request := messages.SetIdentityURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting identity url failed: " + result.(messages.ActionResponse).Message)
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
		request := messages.SetNotificationsURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Done")
			} else {
				println("Setting notifications url failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			println("Wrong IPC response type")
		}

	},
}

var getRuntimeConfigCmd = &cobra.Command{
	Use:   "get-runtime-config",
	Short: "Get the runtime config",
	Long:  `Get the runtime config.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.GetRuntimeConfigRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result := result.(type) {
		case messages.GetRuntimeConfigResponse:
			fmt.Println("{")
			fmt.Println("  \"useMemguard\": " + fmt.Sprintf("%t", result.UseMemguard) + ",")
			fmt.Println("  \"SSHAgentSocketPath\": \"" + result.SSHAgentSocketPath + "\",")
			fmt.Println("  \"goldwardenSocketPath\": \"" + result.GoldwardenSocketPath + "\"")
			fmt.Println("}")
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
	configCmd.AddCommand(getRuntimeConfigCmd)
}
