package cmd

import (
	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage the vault pin",
	Long:  `Manage the vault pin. The pin is used to unlock the vault.`,
}

var setPinCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a new pin",
	Long:  `Set a new pin. The pin is used to unlock the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := commandClient.SendToAgent(messages.UpdateVaultPINRequest{})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Pin updated")
			} else {
				println("Pin updating failed: " + result.(messages.ActionResponse).Message)
			}
		default:
			println("Wrong response type")
		}
	},
}

var pinStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if a pin is set",
	Long:  `Check if a pin is set. The pin is used to unlock the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := commandClient.SendToAgent(messages.GetVaultPINRequest{})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			println("Pin status: " + result.(messages.ActionResponse).Message)
		default:
			println("Wrong response type")
		}
	},
}

func init() {
	vaultCmd.AddCommand(pinCmd)
	pinCmd.AddCommand(setPinCmd)
	pinCmd.AddCommand(pinStatusCmd)
}
