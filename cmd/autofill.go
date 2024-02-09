//go:build linux

package cmd

import (
	"os"

	"github.com/quexten/goldwarden/autotype"
	"github.com/spf13/cobra"
)

var autofillCmd = &cobra.Command{
	Hidden: true,
	Use:    "autotype",
	Short:  "Autotype credentials",
	Long:   `Autotype credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		// get pasword from env
		password := os.Getenv("PASSWORD")
		autotype.TypeString(username + "\t" + password)
	},
}

func init() {
	rootCmd.AddCommand(autofillCmd)
	autofillCmd.PersistentFlags().String("username", "", "")
}
