package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/quexten/goldwarden/cli/ipc/messages"
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
				fmt.Println("Unlocked")
			} else {
				fmt.Println("Not unlocked: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong response type")
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
				fmt.Println("Locked")
			} else {
				fmt.Println("Not locked: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong response type")
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
				fmt.Println("Purged")
			} else {
				fmt.Println("Not purged: " + result.(messages.ActionResponse).Message)
			}
		default:
			fmt.Println("Wrong response type")
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
			response := map[string]interface{}{}
			status := result.(messages.VaultStatusResponse)
			response["locked"] = status.Locked
			response["loginEntries"] = status.NumberOfLogins
			response["noteEntries"] = status.NumberOfNotes
			response["lastSynced"] = time.Unix(status.LastSynced, 0).String()
			response["websocketConnected"] = status.WebsockedConnected
			response["pinSet"] = status.PinSet
			response["loggedIn"] = status.LoggedIn
			responseJSON, _ := json.Marshal(response)
			fmt.Println(string(responseJSON))
		default:
			fmt.Println("Wrong response type")
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
