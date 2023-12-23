//go:build linux || freebsd

package processsecurity

import (
	"fmt"

	"github.com/godbus/dbus/v5"
	"golang.org/x/sys/unix"
)

func DisableDumpable() error {
	return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
}

func MonitorLocks(onlock func()) error {
	bus, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	err = bus.AddMatchSignal(dbus.WithMatchInterface("org.gnome.ScreenSaver"))
	if err != nil {
		return err
	}
	err = bus.AddMatchSignal(dbus.WithMatchMember("org.freedesktop.ScreenSaver"))
	if err != nil {
		return err
	}

	signals := make(chan *dbus.Signal, 10)
	bus.Signal(signals)
	for {
		select {
		case message := <-signals:
			fmt.Println("Message:", message)
			fmt.Println("name ", message.Name)
			if message.Name == "org.gnome.ScreenSaver.ActiveChanged" {
				if len(message.Body) == 0 {
					continue
				}
				locked, err := message.Body[0].(bool)
				if err || locked {
					onlock()
				}
			}
			if message.Name == "org.freedesktop.ScreenSaver.ActiveChanged" {
				if len(message.Body) == 0 {
					continue
				}
				locked, err := message.Body[0].(bool)
				if err || locked {
					onlock()
				}
			}
		}
	}

	return nil
}
