package cmd

import (
	"fmt"
	"syscall"

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

		cleanup := func() {
			fmt.Println("removing sockets and exiting")
			fmt.Println("unlinking", runtimeConfig.GoldwardenSocketPath)
			err := syscall.Unlink(runtimeConfig.GoldwardenSocketPath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("unlinking", runtimeConfig.SSHAgentSocketPath)
			err = syscall.Unlink(runtimeConfig.SSHAgentSocketPath)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("memguard wiping memory and exiting")
			memguard.SafeExit(0)
		}

		err := agent.StartUnixAgent(runtimeConfig.GoldwardenSocketPath, runtimeConfig)
		if err != nil {
			panic(err)
		}
		cleanup()
	},
}

func init() {
	rootCmd.AddCommand(daemonizeCmd)
}
