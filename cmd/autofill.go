//go:build linux

package cmd

import (
	"github.com/quexten/goldwarden/autotype"
	"github.com/spf13/cobra"
)

var autofillCmd = &cobra.Command{
	Use:   "autotype",
	Short: "Autotype credentials",
	Long:  `Autotype credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		autotype.TypeString(username + "\t" + password)
	},
}

func init() {
	rootCmd.AddCommand(autofillCmd)
	autofillCmd.PersistentFlags().String("username", "", "")
	autofillCmd.PersistentFlags().String("password", "", "")
}
