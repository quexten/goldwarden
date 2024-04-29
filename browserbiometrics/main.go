package browserbiometrics

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/browserbiometrics/logging"
)

//go:embed mozilla-com.8bit.bitwarden.json
var templateMozilla string

//go:embed chrome-com.8bit.bitwarden.json
var templateChrome string

//go:embed goldwarden-proxy.sh
var proxyScript string

var chromiumPaths = []string{
	"~/.config/google-chrome/",
	"~/.config/google-chrome-beta/",
	"~/.config/google-chrome-unstable/",
	"~/.config/chromium/",
	"~/.config/BraveSoftware/Brave-Browser/",
	"~/.config/thorium/",
	"~/.config/microsoft-edge-beta/",
	"~/.config/microsoft-edge-dev/",
}
var mozillaPaths = []string{"~/.mozilla/", "~/.librewolf/", "~/.waterfox/"}

const appID = "com.quexten.bw-bio-handler"

var transportKey []byte

func Main(rtCfg *config.RuntimeConfig) error {
	logging.Debugf("Starting browserbiometrics")
	var err error
	transportKey, err = generateTransportKey()
	if err != nil {
		return err
	}
	logging.Debugf("Generated transport key")

	setupCommunication()
	return readLoop(rtCfg)
}

func DetectAndInstallBrowsers() error {
	var err error

	// first, ensure the native messaging hosts dirs exist
	for _, path := range chromiumPaths {
		path = strings.ReplaceAll(path, "~", os.Getenv("HOME"))
		_, err = os.Stat(path)
		if err != nil {
			continue
		}

		_, err = os.Stat(path + "NativeMessagingHosts/")
		if err == nil {
			fmt.Println("Native messaging host directory already exists: " + path + "NativeMessagingHosts/")
			continue
		}
		err = os.MkdirAll(path+"NativeMessagingHosts/", 0755)
		if err != nil {
			fmt.Println("Error creating native messaging host directory: " + err.Error())
		} else {
			fmt.Println("Created native messaging host directory: " + path + "NativeMessagingHosts/")
		}
	}
	for _, path := range mozillaPaths {
		path = strings.ReplaceAll(path, "~", os.Getenv("HOME"))
		_, err = os.Stat(path)
		if err != nil {
			continue
		}

		_, err = os.Stat(path + "native-messaging-hosts/")
		if err == nil {
			fmt.Println("Native messaging host directory already exists: " + path + "native-messaging-hosts/")
			continue
		}
		err = os.MkdirAll(path+"native-messaging-hosts/", 0755)
		if err != nil {
			fmt.Println("Error creating native messaging host directory: " + err.Error())
		} else {
			fmt.Println("Created native messaging host directory: " + path + "native-messaging-hosts/")
		}
	}

	err = detectAndInstallBrowsers(".config")
	if err != nil {
		return err
	}
	for _, path := range mozillaPaths {
		path = strings.ReplaceAll(path, "~/", "")
		err = detectAndInstallBrowsers(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func detectAndInstallBrowsers(startPath string) error {
	home := os.Getenv("HOME")
	err := filepath.Walk(home+"/"+startPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		var tempPath string
		if !strings.HasPrefix(path, home) {
			return nil
		} else {
			tempPath = strings.TrimPrefix(path, home)
		}
		if strings.Count(tempPath, "/") > 5 {
			return nil
		}

		if info.IsDir() && info.Name() == "native-messaging-hosts" {
			fmt.Printf("Found mozilla-like browser: %s\n", path)

			fmt.Println("Removing old manifest and proxy script")
			if _, err := os.Stat(path + "/com.8bit.bitwarden.json"); err == nil {
				err = os.Chmod(path+"/com.8bit.bitwarden.json", 755)
				if err != nil {
					return err
				}
				err = os.Remove(path + "/com.8bit.bitwarden.json")
				if err != nil {
					return err
				}
			}

			if _, err := os.Stat(path + "/goldwarden-proxy.sh"); err == nil {
				err = os.Chmod(path+"/goldwarden-proxy.sh", 755)
				if err != nil {
					return err
				}
				err = os.Remove(path + "/goldwarden-proxy.sh")
				if err != nil {
					return err
				}
			}

			fmt.Println("Writing new manifest")
			manifest := strings.Replace(templateMozilla, "@PATH@", path+"/goldwarden-proxy.sh", 1)
			err = os.WriteFile(path+"/com.8bit.bitwarden.json", []byte(manifest), 0444)
			if err != nil {
				return err
			}

			fmt.Println("Writing new proxy script")
			err = os.WriteFile(path+"/goldwarden-proxy.sh", []byte(proxyScript), 0755)
			if err != nil {
				return err
			}
		} else if info.IsDir() && info.Name() == "NativeMessagingHosts" {
			fmt.Printf("Found chrome-like browser: %s\n", path)

			fmt.Println("Removing old manifest and proxy script")
			if _, err := os.Stat(path + "/com.8bit.bitwarden.json"); err == nil {
				err = os.Chmod(path+"/com.8bit.bitwarden.json", 755)
				if err != nil {
					return err
				}
				err = os.Remove(path + "/com.8bit.bitwarden.json")
				if err != nil {
					return err
				}
			}

			if _, err := os.Stat(path + "/goldwarden-proxy.sh"); err == nil {
				err = os.Chmod(path+"/goldwarden-proxy.sh", 755)
				if err != nil {
					return err
				}
				err = os.Remove(path + "/goldwarden-proxy.sh")
				if err != nil {
					return err
				}
			}

			fmt.Println("Writing new manifest")
			manifest := strings.Replace(templateChrome, "@PATH@", path+"/goldwarden-proxy.sh", 1)
			err = os.WriteFile(path+"/com.8bit.bitwarden.json", []byte(manifest), 0444)
			if err != nil {
				return err
			}

			fmt.Println("Writing new proxy script")
			err = os.WriteFile(path+"/goldwarden-proxy.sh", []byte(proxyScript), 0755)
			if err != nil {
				return err
			}
		}

		return err
	})

	return err
}
