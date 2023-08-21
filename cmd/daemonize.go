package cmd

import (
	"os"
	"os/signal"

	"github.com/awnumar/memguard"
	"github.com/quexten/goldwarden/agent"
	"github.com/spf13/cobra"
)

var daemonizeCmd = &cobra.Command{
	Use:   "daemonize",
	Short: "Starts the agent as a daemon",
	Long: `Starts the agent as a daemon. The agent will run in the background and will
	run in the background until it is stopped.`,
	Run: func(cmd *cobra.Command, args []string) {
		websocketDisabled := runtimeConfig.WebsocketDisabled
		sshDisabled := runtimeConfig.DisableSSHAgent

		if websocketDisabled {
			println("Websocket disabled")
		}

		if sshDisabled {
			println("SSH agent disabled")
		}

		go func() {
			signalChannel := make(chan os.Signal, 1)
			signal.Notify(signalChannel, os.Interrupt)
			<-signalChannel
			memguard.SafeExit(0)
		}()
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		err = agent.StartUnixAgent(home+"/.goldwarden.sock", runtimeConfig)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonizeCmd)
}
