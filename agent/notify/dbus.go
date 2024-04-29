//go:build freebsd || linux

package notify

import (
	"time"

	"github.com/godbus/dbus/v5"
)

var closeListenerMap = make(map[uint32]func())

func notifyDBus(title string, body string, actionName string, timeout time.Duration, onclose func()) error {
	bus, err := dbus.SessionBus()
	if err != nil {
		log.Error("could not get a dbus session: %s", err.Error())
		return err
	}
	obj := bus.Object("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	actions := []string{}
	if actionName != "" {
		actions = append(actions, actionName)
	}

	notificationID++

	call := obj.Call("org.freedesktop.Notifications.Notify", 0, "goldwarden", uint32(notificationID), "", title, body, actions, map[string]dbus.Variant{}, int32(60000))
	if call.Err != nil {
		log.Error("could not call dbus object: %s", call.Err.Error())
		return err
	}
	if len(call.Body) < 1 {
		return err
	}
	id := call.Body[0].(uint32)
	closeListenerMap[id] = onclose

	if timeout == 0 {
		return err
	} else {
		go func(id uint32) {
			time.Sleep(timeout)
			call := obj.Call("org.freedesktop.Notifications.CloseNotification", 0, uint32(id))
			if call.Err != nil {
				return
			}
		}(id)
	}

	return nil
}

func listenForNotificationsDBus() error {
	bus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	err = bus.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.Notifications"))
	if err != nil {
		return err
	}

	signals := make(chan *dbus.Signal, 10)
	bus.Signal(signals)
	for {
		message := <-signals
		if message.Name == "org.freedesktop.Notifications.NotificationClosed" {
			if len(message.Body) < 1 {
				continue
			}
			id, ok := message.Body[0].(uint32)
			if !ok {
				continue
			}
			if id == 0 {
				continue
			}
			if closeListener, ok := closeListenerMap[id]; ok {
				delete(closeListenerMap, id)
				closeListener()
			}
		}
	}
}
