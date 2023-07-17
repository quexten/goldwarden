package cmd

import (
	"github.com/quexten/goldwarden/client"
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

		result, err := client.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
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

		result, err := client.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
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

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration",
	Long:  `Manage the configuration.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(setApiUrlCmd)
	configCmd.AddCommand(setIdentityURLCmd)
}
