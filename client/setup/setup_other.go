//go:build !linux

package setup

import "github.com/quexten/goldwarden/agent/config"

func VerifySetup(runtimeConfig config.RuntimeConfig) bool {
	return true
}
