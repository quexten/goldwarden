package cmd

import (
	"fmt"

	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Manage the vault",
	Long:  `Manage the vault.`,
}

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlocks the vault",
	Long:  `Unlocks the vault. You will be prompted for your pin. The pin is empty by default.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.UnlockVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Unlocked")
			} else {
				println("Not unlocked: " + result.(messages.ActionResponse).Message)
			}
		default:
			println("Wrong response type")
		}
	},
}

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Locks the vault",
	Long:  `Locks the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.LockVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Locked")
			} else {
				println("Not locked: " + result.(messages.ActionResponse).Message)
			}
		default:
			println("Wrong response type")
		}
	},
}

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Wipes the vault",
	Long:  `Wipes the vault and encryption keys from ram and config. Does not delete any entries on the server side.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.WipeVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ActionResponse:
			if result.(messages.ActionResponse).Success {
				println("Purged")
			} else {
				println("Not purged: " + result.(messages.ActionResponse).Message)
			}
		default:
			println("Wrong response type")
		}
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Shows the vault status",
	Long:  `Shows the vault status.`,
	Run: func(cmd *cobra.Command, args []string) {
		request := messages.VaultStatusRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.VaultStatusResponse:
			status := result.(messages.VaultStatusResponse)
			fmt.Println("Locked: ", status.Locked)
			fmt.Println("Number of logins: ", status.NumberOfLogins)
			fmt.Println("Number of notes: ", status.NumberOfNotes)
		default:
			println("Wrong response type")
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)
	vaultCmd.AddCommand(unlockCmd)
	vaultCmd.AddCommand(lockCmd)
	vaultCmd.AddCommand(purgeCmd)
	vaultCmd.AddCommand(statusCmd)
}
