package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quexten/goldwarden/cli/ipc/messages"
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
				fmt.Println("Done")
			} else {
				fmt.Println("Setting api url failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
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
				fmt.Println("Done")
			} else {
				fmt.Println("Setting identity url failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
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
				fmt.Println("Done")
			} else {
				fmt.Println("Setting notifications url failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
		}

	},
}

var setVaultURLCmd = &cobra.Command{
	Use:   "set-vault-url",
	Short: "Set the vault url",
	Long:  `Set the vault url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		url := args[0]
		request := messages.SetVaultURLRequest{}
		request.Value = url

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				fmt.Println("Done")
			} else {
				fmt.Println("Setting vault url failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
		}

	},
}

var setURLsAutomaticallyCmd = &cobra.Command{
	Use:   "set-server",
	Short: "Set the urls automatically",
	Long:  `Set the api/identity/vault/notification urls automatically from a base url.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		value := args[0]
		request := messages.SetURLsAutomaticallyRequest{}
		request.Value = value

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				fmt.Println("Done")
			} else {
				fmt.Println("Setting urls automatically failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
		}

	},
}

var getEnvironmentCmd = &cobra.Command{
	Use:   "get-environment",
	Short: "Get the environment",
	Long:  `Get the environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.GetConfigEnvironmentRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result := result.(type) {
		case messages.GetConfigEnvironmentResponse:
			response := map[string]string{}
			response["api"] = result.ApiURL
			response["identity"] = result.IdentityURL
			response["notifications"] = result.NotificationsURL
			response["vault"] = result.VaultURL
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
		default:
			fmt.Println("Wrong IPC response type")
		}
	},
}

var setApiClientIDCmd = &cobra.Command{
	Use:   "set-client-id",
	Short: "Set the client id",
	Long:  `Set the client id.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		id := args[0]
		if len(id) >= 2 && strings.HasPrefix(id, "\"") && strings.HasSuffix(id, "\"") {
			id = id[1 : len(id)-1]
		}
		id = strings.TrimSpace(id)
		request := messages.SetClientIDRequest{}
		request.Value = id

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				fmt.Println("Done")
			} else {
				fmt.Println("Setting api client id failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
		}

	},
}

var setApiSecretCmd = &cobra.Command{
	Use:   "set-client-secret",
	Short: "Set the api secret",
	Long:  `Set the api secret.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		secret := args[0]
		if len(secret) >= 2 && strings.HasPrefix(secret, "\"") && strings.HasSuffix(secret, "\"") {
			secret = secret[1 : len(secret)-1]
		}
		secret = strings.TrimSpace(secret)
		request := messages.SetClientSecretRequest{}
		request.Value = secret

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				fmt.Println("Done")
			} else {
				fmt.Println("Setting api secret failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong IPC response type")
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
			response := map[string]interface{}{}
			response["useMemguard"] = result.UseMemguard
			response["SSHAgentSocketPath"] = result.SSHAgentSocketPath
			response["goldwardenSocketPath"] = result.GoldwardenSocketPath
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
		default:
			fmt.Println("Wrong IPC response type")
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
	configCmd.AddCommand(setVaultURLCmd)
	configCmd.AddCommand(setURLsAutomaticallyCmd)
	configCmd.AddCommand(getEnvironmentCmd)
	configCmd.AddCommand(getRuntimeConfigCmd)
	configCmd.AddCommand(setApiClientIDCmd)
	configCmd.AddCommand(setApiSecretCmd)
}
