/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/quexten/goldwarden/ipc/messages"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a command with environment variables from your vault",
	Long: `Runs a command with environment variables from your vault.
	The variables are stored as a secure note. Consult the documentation for more information.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Error: No command specified")
			return
		}

		executable := args[0]
		executableArgs := args[1:]

		env := []string{}

		result, err := commandClient.SendToAgent(messages.GetCLICredentialsRequest{
			ApplicationName: executable,
		})
		if err != nil {
			handleSendToAgentError(err)
			return
		}

		switch result.(type) {
		case messages.GetCLICredentialsResponse:
			response := result.(messages.GetCLICredentialsResponse)
			for key, value := range response.Env {
				env = append(env, key+"="+value)
			}
		case messages.ActionResponse:
			fmt.Println("Error: " + result.(messages.ActionResponse).Message)
			return
		}

		command := exec.Command(executable, executableArgs...)
		command.Env = append(command.Env, os.Environ()...)
		command.Env = append(command.Env, env...)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		command.Stdin = os.Stdin
		_ = command.Run()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
