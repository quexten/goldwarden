package cmd

import (
	"fmt"

	_ "embed"

	"github.com/spf13/cobra"
)

//go:embed version.txt
var version string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Shows the version of the cli",
	Long:  `Shows the version of the cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
