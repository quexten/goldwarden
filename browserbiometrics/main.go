package browserbiometrics

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quexten/goldwarden/browserbiometrics/logging"
)

const appID = "com.quexten.bw-bio-handler"

var transportKey []byte

func Main() {
	if os.Args[1] == "install" {
		var err error
		err = detectAndInstallBrowsers(".config")
		if err != nil {
			panic("Failed to detect browsers: " + err.Error())
		}
		err = detectAndInstallBrowsers(".mozilla")
		if err != nil {
			panic("Failed to detect browsers: " + err.Error())
		}
		return
	}

	logging.Debugf("Starting browserbiometrics")
	transportKey = generateTransportKey()
	logging.Debugf("Generated transport key")

	setupCommunication()
	readLoop()
}

func DetectAndInstallBrowsers() error {
	var err error
	err = detectAndInstallBrowsers(".config")
	if err != nil {
		return err
	}
	err = detectAndInstallBrowsers(".mozilla")
	if err != nil {
		return err
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
		if strings.Count(tempPath, "/") > 3 {
			return nil
		}

		binPath, err := os.Executable()
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == "native-messaging-hosts" {
			fmt.Printf("Found mozilla-like browser: %s\n", path)
			manifest := strings.Replace(templateMozilla, "PATH", binPath, 1)
			err = os.WriteFile(path+"/com.8bit.bitwarden.json", []byte(manifest), 0644)
		} else if info.IsDir() && info.Name() == "NativeMessagingHosts" {
			fmt.Printf("Found chrome-like browser: %s\n", path)
			manifest := strings.Replace(templateChrome, "PATH", binPath, 1)
			err = os.WriteFile(path+"/com.8bit.bitwarden.json", []byte(manifest), 0644)
		}

		return err
	})

	return err
}
