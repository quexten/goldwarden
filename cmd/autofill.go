//go:build autofill

package cmd

import (
	"github.com/quexten/goldwarden/autofill"
	"github.com/spf13/cobra"
)

var autofillCmd = &cobra.Command{
	Use:   "autofill",
	Short: "Autofill credentials",
	Long:  `Autofill credentials`,
	Run: func(cmd *cobra.Command, args []string) {
		layout := cmd.Flag("layout").Value.String()
		autofill.Run(layout)
	},
}

func init() {
	rootCmd.AddCommand(autofillCmd)
	autofillCmd.PersistentFlags().String("layout", "qwerty", "")
}
