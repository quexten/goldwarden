package cmd

import (
	"github.com/quexten/goldwarden/ipc"
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
		request := ipc.UnlockVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Unlocked")
			} else {
				println("Not unlocked: " + result.(ipc.ActionResponse).Message)
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
		request := ipc.LockVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Locked")
			} else {
				println("Not locked: " + result.(ipc.ActionResponse).Message)
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
		request := ipc.WipeVaultRequest{}

		result, err := commandClient.SendToAgent(request)
		if err != nil {
			println("Error: " + err.Error())
			println("Is the daemon running?")
			return
		}

		switch result.(type) {
		case ipc.ActionResponse:
			if result.(ipc.ActionResponse).Success {
				println("Purged")
			} else {
				println("Not purged: " + result.(ipc.ActionResponse).Message)
			}
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
}
