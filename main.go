package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/browserbiometrics"
	"github.com/quexten/goldwarden/cmd"
)

func main() {
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
		UseMemguard:           os.Getenv("GOLDWARDEN_NO_MEMGUARD") != "true",
		SSHAgentSocketPath:    os.Getenv("GOLDWARDEN_SSH_AUTH_SOCK"),
		GoldwardenSocketPath:  os.Getenv("GOLDWARDEN_SOCKET_PATH"),
		DaemonAuthToken:       os.Getenv("GOLDWARDEN_DAEMON_AUTH_TOKEN"),

		ConfigDirectory: configPath,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	if runtimeConfig.SSHAgentSocketPath == "" {
		runtimeConfig.SSHAgentSocketPath = home + "/.goldwarden-ssh-agent.sock"
	}
	if runtimeConfig.GoldwardenSocketPath == "" {
		runtimeConfig.GoldwardenSocketPath = home + "/.goldwarden.sock"
	}

	if len(os.Args) > 1 && (strings.Contains(os.Args[1], "com.8bit.bitwarden.json") || strings.Contains(os.Args[1], "chrome-extension://")) {
		browserbiometrics.Main(&runtimeConfig)
		return
	}

	_, err = os.Stat("/.flatpak-info")
	isFlatpak := err == nil
	if isFlatpak {
		userHome, _ := os.UserHomeDir()
		runtimeConfig.ConfigDirectory = userHome + "/.var/app/com.quexten.Goldwarden/config/goldwarden.json"
		runtimeConfig.ConfigDirectory = strings.ReplaceAll(runtimeConfig.ConfigDirectory, "~", userHome)
		fmt.Println("Flatpak Config directory: " + runtimeConfig.ConfigDirectory)
		runtimeConfig.SSHAgentSocketPath = userHome + "/.var/app/com.quexten.Goldwarden/data/ssh-auth-sock"
		runtimeConfig.GoldwardenSocketPath = userHome + "/.var/app/com.quexten.Goldwarden/data/goldwarden.sock"
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

	cmd.Execute(runtimeConfig)
}
