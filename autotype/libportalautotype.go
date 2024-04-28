//go:build linux

package autotype

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/quexten/goldwarden/logging"
)

var globalID = 0

const autoTypeDelay = 1 * time.Millisecond

var log = logging.GetLogger("Goldwarden", "Autotype")

// todo need to store this encrypted. will be done when migrating this file to python
func persistToken(token string) error {
	tokenPath := ""
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat("/.flatpak-info"); err == nil {
		tokenPath = userHome + "/.var/app/com.quexten.Goldwarden/config/autotype_token.txt"
	} else {
		tokenPath = userHome + "/.config/goldwarden/autotype_token.txt"
	}

	err = ioutil.WriteFile(tokenPath, []byte(token), 0644)
	if err != nil {
		return err
	}
	return nil
}

func readToken() (string, error) {
	tokenPath := ""
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat("/.flatpak-info"); err == nil {
		tokenPath = userHome + "/.var/app/com.quexten.Goldwarden/config/autotype_token.txt"
	} else {
		tokenPath = userHome + "/.config/goldwarden/autotype_token.txt"
	}

	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

func TypeString(textToType string) {
	log.Info("Starting to Type String")
	bus, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	obj := bus.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	obj.AddMatchSignal("org.freedesktop.portal.Request", "Response")

	globalID++
	res0 := obj.Call("org.freedesktop.portal.RemoteDesktop.CreateSession", 0, map[string]dbus.Variant{
		"session_handle_token": dbus.MakeVariant("u" + fmt.Sprint(globalID)),
	})
	if res0.Err != nil {
		log.Error("Error creating session: %s", res0.Err.Error())
		return
	}

	signals := make(chan *dbus.Signal, 10)
	bus.Signal(signals)

	var state = 0
	var sessionHandle dbus.ObjectPath

	for {
		message := <-signals
		switch state {
		case 0:
			log.Info("Selecting Devices")
			result := message.Body[1].(map[string]dbus.Variant)
			resultSessionHandle := result["session_handle"]
			sessionHandle = dbus.ObjectPath(resultSessionHandle.String()[1 : len(resultSessionHandle.String())-1])
			options := map[string]dbus.Variant{
				"types":        dbus.MakeVariant(uint32(1)),
				"persist_mode": dbus.MakeVariant(uint32(2)),
			}
			if token, err := readToken(); err == nil {
				log.Info("Restoring token, no confirmation prompt")
				options["restore_token"] = dbus.MakeVariant(token)
			}

			res := obj.Call("org.freedesktop.portal.RemoteDesktop.SelectDevices", 0, sessionHandle, options)
			if res.Err != nil {
				log.Error("Error selecting devices: %s", res.Err.Error())
			}
			state = 1
		case 1:
			log.Info("Starting Session")
			res := obj.Call("org.freedesktop.portal.RemoteDesktop.Start", 0, sessionHandle, "", map[string]dbus.Variant{})
			if res.Err != nil {
				log.Error("Error starting session: %s", res.Err.Error())
			}
			state = 2
		case 2:
			// try to cast to interface array
			if len(message.Body) == 2 {
				if resMap, ok := message.Body[1].(map[string]dbus.Variant); ok {
					// check if restore token in response
					if restoreToken, ok := resMap["restore_token"]; ok {
						token := restoreToken.Value().(string)
						if err := persistToken(token); err != nil {
							log.Error("Error persisting token: %s", err.Error())
						}
					}
				}
			}

			log.Info("Performing Typing")
			time.Sleep(1000 * time.Millisecond)
			for _, char := range textToType {
				if char == '\t' {
					obj.Call("org.freedesktop.portal.RemoteDesktop.NotifyKeyboardKeycode", 0, sessionHandle, map[string]dbus.Variant{}, 15, uint32(1))
					time.Sleep(autoTypeDelay)
					obj.Call("org.freedesktop.portal.RemoteDesktop.NotifyKeyboardKeycode", 0, sessionHandle, map[string]dbus.Variant{}, 15, uint32(0))
					time.Sleep(autoTypeDelay)
				} else {
					obj.Call("org.freedesktop.portal.RemoteDesktop.NotifyKeyboardKeysym", 0, sessionHandle, map[string]dbus.Variant{}, int32(char), uint32(1))
					time.Sleep(autoTypeDelay)
					obj.Call("org.freedesktop.portal.RemoteDesktop.NotifyKeyboardKeysym", 0, sessionHandle, map[string]dbus.Variant{}, int32(char), uint32(0))
					time.Sleep(autoTypeDelay)
				}
			}
			bus.Close()
			return
		default:
			return
		}
	}
}
