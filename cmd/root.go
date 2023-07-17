package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

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
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
