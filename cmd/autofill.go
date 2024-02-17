//go:build linux

package cmd

import (
	"bufio"
	"encoding/hex"
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
		reader := bufio.NewReader(os.Stdin)
		textHex, _ := reader.ReadString('\n')
		text, _ := hex.DecodeString(textHex)
		autotype.TypeString(string(text))
	},
}

func init() {
	rootCmd.AddCommand(autofillCmd)
}
