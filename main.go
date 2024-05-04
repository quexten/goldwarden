package main

import (
	"os"
	"strings"

	"github.com/quexten/goldwarden/cli/agent/config"
	"github.com/quexten/goldwarden/cli/browserbiometrics"
	"github.com/quexten/goldwarden/cli/cmd"
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
		WebsocketDisabled:    os.Getenv("GOLDWARDEN_WEBSOCKET_DISABLED") == "true",
		DisableSSHAgent:      os.Getenv("GOLDWARDEN_SSH_AGENT_DISABLED") == "true",
		DoNotPersistConfig:   os.Getenv("GOLDWARDEN_DO_NOT_PERSIST_CONFIG") == "true",
		DeviceUUID:           os.Getenv("GOLDWARDEN_DEVICE_UUID"),
		AuthMethod:           os.Getenv("GOLDWARDEN_AUTH_METHOD"),
		User:                 os.Getenv("GOLDWARDEN_AUTH_USER"),
		Password:             os.Getenv("GOLDWARDEN_AUTH_PASSWORD"),
		Pin:                  os.Getenv("GOLDWARDEN_PIN"),
		UseMemguard:          os.Getenv("GOLDWARDEN_NO_MEMGUARD") != "true",
		SSHAgentSocketPath:   os.Getenv("GOLDWARDEN_SSH_AUTH_SOCK"),
		GoldwardenSocketPath: os.Getenv("GOLDWARDEN_SOCKET_PATH"),
		DaemonAuthToken:      os.Getenv("GOLDWARDEN_DAEMON_AUTH_TOKEN"),

		ConfigDirectory: configPath,
	}

	_, err := os.Stat("/.flatpak-info")
	isFlatpak := err == nil
	if isFlatpak {
		userHome, _ := os.UserHomeDir()
		runtimeConfig.ConfigDirectory = userHome + "/.var/app/com.quexten.Goldwarden/config/goldwarden.json"
		runtimeConfig.ConfigDirectory = strings.ReplaceAll(runtimeConfig.ConfigDirectory, "~", userHome)
	}

	if len(os.Args) > 1 && (strings.Contains(os.Args[1], "com.8bit.bitwarden.json") || strings.Contains(os.Args[1], "chrome-extension://")) {
		err = browserbiometrics.Main(&runtimeConfig)
		if err != nil {
			panic(err)
		}
	}

	cmd.Execute(runtimeConfig)
}
