//go:build linux

package autotype

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/quexten/goldwarden/logging"
)

var globalID = 0

const autoTypeDelay = 1 * time.Millisecond

var log = logging.GetLogger("Goldwarden", "Autotype")

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
			res := obj.Call("org.freedesktop.portal.RemoteDesktop.SelectDevices", 0, sessionHandle, map[string]dbus.Variant{
				"types": dbus.MakeVariant(uint32(1)),
			})
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
