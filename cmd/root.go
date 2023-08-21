package cmd

import (
	"os"

	"github.com/quexten/goldwarden/agent"
	"github.com/quexten/goldwarden/client"
	"github.com/spf13/cobra"
)

var commandClient client.Client

var rootCmd = &cobra.Command{
	Use:   "goldwarden",
	Short: "OS level integration for Bitwarden",
	Long: `Goldwarden is a daemon that runs in the background and provides
	OS level integration for Bitwarden, such as SSH agent integration, 
	biometric unlock, and more.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	goldwardenSingleProcess := os.Getenv("GOLDWARDEN_SINGLE_PROCESS")
	if goldwardenSingleProcess == "true" {
		recv, send := agent.StartVirtualAgent()
		commandClient = client.NewVirtualClient(send, recv)
	} else {
		commandClient = client.NewUnixSocketClient()
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
