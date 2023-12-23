//go:build linux

package autotype

import (
	"fmt"
	"time"

	"github.com/godbus/dbus/v5"
)

var globalID = 0

const autoTypeDelay = 1 * time.Millisecond

func TypeString(textToType string) {
	bus, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	obj := bus.Object("org.freedesktop.portal.Desktop", "/org/freedesktop/portal/desktop")
	obj.AddMatchSignal("org.freedesktop.portal.Request", "Response")

	globalID++
	obj.Call("org.freedesktop.portal.RemoteDesktop.CreateSession", 0, map[string]dbus.Variant{
		"session_handle_token": dbus.MakeVariant("u" + fmt.Sprint(globalID)),
	})

	signals := make(chan *dbus.Signal, 10)
	bus.Signal(signals)

	var state = 0
	var sessionHandle dbus.ObjectPath

	for {
		select {
		case message := <-signals:
			fmt.Println("Message:", message)
			if state == 0 {
				result := message.Body[1].(map[string]dbus.Variant)
				resultSessionHandle := result["session_handle"]
				sessionHandle = dbus.ObjectPath(resultSessionHandle.String()[1 : len(resultSessionHandle.String())-1])
				obj.Call("org.freedesktop.portal.RemoteDesktop.SelectDevices", 0, sessionHandle, map[string]dbus.Variant{})
				state = 1
			} else if state == 1 {
				obj.Call("org.freedesktop.portal.RemoteDesktop.Start", 0, sessionHandle, "", map[string]dbus.Variant{})
				state = 2
			} else if state == 2 {
				state = 3
				time.Sleep(200 * time.Millisecond)
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
			}
		}
	}
}
