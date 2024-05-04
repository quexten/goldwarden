//go:build linux

package setup

import (
	"fmt"

	"github.com/quexten/goldwarden/cli/agent/config"
	"github.com/quexten/goldwarden/cli/cmd"
)

func VerifySetup(runtimeConfig config.RuntimeConfig) bool {
	if !cmd.IsPolkitSetup() {
		fmt.Println("Polkit is not setup. Run 'goldwarden setup polkit' to set it up.")
		return false
	}

	return true
}
