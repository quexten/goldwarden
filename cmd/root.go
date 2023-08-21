package cmd

import (
	"os"

	"github.com/quexten/goldwarden/agent"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc"
	"github.com/spf13/cobra"
)

var commandClient client.Client
var runtimeConfig config.RuntimeConfig

var rootCmd = &cobra.Command{
	Use:   "goldwarden",
	Short: "OS level integration for Bitwarden",
	Long: `Goldwarden is a daemon that runs in the background and provides
	OS level integration for Bitwarden, such as SSH agent integration, 
	biometric unlock, and more.`,
}

func Execute(cfg config.RuntimeConfig) {
	runtimeConfig = cfg

	goldwardenSingleProcess := os.Getenv("GOLDWARDEN_SINGLE_PROCESS")
	if goldwardenSingleProcess == "true" {
		recv, send := agent.StartVirtualAgent(runtimeConfig)
		commandClient = client.NewVirtualClient(send, recv)
	} else {
		commandClient = client.NewUnixSocketClient()
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func loginIfRequired() error {
	var err error

	if runtimeConfig.AuthMethod == "password" {
		_, err = commandClient.SendToAgent(ipc.DoLoginRequest{
			Email:    runtimeConfig.User,
			Password: runtimeConfig.Password,
		})
	} else if runtimeConfig.AuthMethod == "passwordless" {
		_, err = commandClient.SendToAgent(ipc.DoLoginRequest{
			Email:        runtimeConfig.User,
			Passwordless: true,
		})
	}

	return err
}
