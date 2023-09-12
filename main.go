package main

import (
	"os"
	"strings"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/browserbiometrics"
	"github.com/quexten/goldwarden/client/setup"
	"github.com/quexten/goldwarden/cmd"
)

func main() {
	if len(os.Args) > 1 && strings.Contains(os.Args[1], "com.8bit.bitwarden.json") {
		browserbiometrics.Main()
		return
	}

	var configPath string
	if path, found := os.LookupEnv("GOLDWARDEN_CONFIG_DIRECTORY"); found {
		configPath = path
	} else {
		configPath = config.DefaultConfigPath
	}
	userHome, _ := os.UserHomeDir()
	configPath = strings.ReplaceAll(configPath, "~", userHome)

	runtimeConfig := config.RuntimeConfig{
		WebsocketDisabled:     os.Getenv("GOLDWARDEN_WEBSOCKET_DISABLED") == "true",
		DisableSSHAgent:       os.Getenv("GOLDWARDEN_SSH_AGENT_DISABLED") == "true",
		DisableAuth:           os.Getenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED") == "true",
		DisablePinRequirement: os.Getenv("GOLDWARDEN_PIN_REQUIREMENT_DISABLED") == "true",
		DoNotPersistConfig:    os.Getenv("GOLDWARDEN_DO_NOT_PERSIST_CONFIG") == "true",
		ApiURI:                os.Getenv("GOLDWARDEN_API_URI"),
		IdentityURI:           os.Getenv("GOLDWARDEN_IDENTITY_URI"),
		SingleProcess:         os.Getenv("GOLDWARDEN_SINGLE_PROCESS") == "true",
		DeviceUUID:            os.Getenv("GOLDWARDEN_DEVICE_UUID"),
		AuthMethod:            os.Getenv("GOLDWARDEN_AUTH_METHOD"),
		User:                  os.Getenv("GOLDWARDEN_AUTH_USER"),
		Password:              os.Getenv("GOLDWARDEN_AUTH_PASSWORD"),
		Pin:                   os.Getenv("GOLDWARDEN_PIN"),

		ConfigDirectory: configPath,
	}

	if runtimeConfig.SingleProcess {
		runtimeConfig.DisablePinRequirement = true
		runtimeConfig.DisableAuth = true
	}

	if runtimeConfig.DisablePinRequirement {
		runtimeConfig.DoNotPersistConfig = true
	}

	if runtimeConfig.DisableAuth {
		os.Setenv("GOLDWARDEN_SYSTEM_AUTH_DISABLED", "true")
	}

	if !setup.VerifySetup(runtimeConfig) {
		return
	}

	cmd.Execute(runtimeConfig)
}
