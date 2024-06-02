package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/quexten/goldwarden/cli/ipc/messages"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Commands for managing SSH keys",
	Long:  `Commands for managing SSH keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// runCmd represents the run command
var sshAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Runs a command with environment variables from your vault",
	Long: `Runs a command with environment variables from your vault.
	The variables are stored as a secure note. Consult the documentation for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := loginIfRequired()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		name, _ := cmd.Flags().GetString("name")
		copyToClipboard, _ := cmd.Flags().GetBool("clipboard")

		result, err := commandClient.SendToAgent(messages.CreateSSHKeyRequest{
			Name: name,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.CreateSSHKeyResponse:
			response := result.(messages.CreateSSHKeyResponse)
			fmt.Println(response.Digest)

			if copyToClipboard {
				err := clipboard.WriteAll(string(response.Digest))
				if err != nil {
					panic(err)
				}
			}
			return
		case messages.ActionResponse:
			fmt.Println("Error: " + result.(messages.ActionResponse).Message)
			return
		}
	},
}

var listSSHCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all SSH keys in your vault",
	Long:  `Lists all SSH keys in your vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := loginIfRequired()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		result, err := commandClient.SendToAgent(messages.GetSSHKeysRequest{})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.GetSSHKeysResponse:
			response := result.(messages.GetSSHKeysResponse)
			for _, key := range response.Keys {
				fmt.Println(key)
			}
			return
		case messages.ActionResponse:
			fmt.Println("Error: " + result.(messages.ActionResponse).Message)
			return
		}
	},
}

var importSSHCmd = &cobra.Command{
	Use:   "import",
	Short: "Imports an SSH key into your vault",
	Long:  `Imports an SSH key into your vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]
		fmt.Println("Importing SSH key from " + filename)

		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			name = "Imported SSH Key"
		}

		if _, err := os.Stat(filename); os.IsNotExist(err) {
			fmt.Println("Error: File does not exist")
			return
		}

		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(file)
		key := buf.String()

		result, err := commandClient.SendToAgent(messages.ImportSSHKeyRequest{
			Key:  key,
			Name: name,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.ImportSSHKeyResponse:
			response := result.(messages.ImportSSHKeyResponse)
			if response.Success {
				fmt.Println("Success")
			} else {
				fmt.Println("Error: " + response.ErrorMsg)
			}
			return
		case messages.ActionResponse:
			fmt.Println("Error: " + result.(messages.ActionResponse).Message)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
	sshCmd.AddCommand(sshAddCmd)
	sshAddCmd.PersistentFlags().String("name", "", "")
	_ = sshAddCmd.MarkFlagRequired("name")
	sshAddCmd.PersistentFlags().Bool("clipboard", false, "Copy the public key to the clipboard")
	sshCmd.AddCommand(listSSHCmd)
	importSSHCmd.PersistentFlags().String("name", "", "")
	sshCmd.AddCommand(importSSHCmd)
}
