//go:build linux || freebsd

package processsecurity

import (
	"time"

	"github.com/godbus/dbus/v5"
	"github.com/quexten/goldwarden/cli/agent/processsecurity/isdelve"
	"golang.org/x/sys/unix"
)

const IDLE_TIME = 60 * 15

func DisableDumpable() error {
	if isdelve.Enabled {
		return nil
	} else {
		return unix.Prctl(unix.PR_SET_DUMPABLE, 0, 0, 0, 0)
	}
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
	err = bus.AddMatchSignal(dbus.WithMatchInterface("org.freedesktop.ScreenSaver"))
	if err != nil {
		return err
	}

	signals := make(chan *dbus.Signal, 10)
	bus.Signal(signals)
	for {
		message := <-signals
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

func MonitorIdle(onidle func()) error {
	bus, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	var wasidle = false
	for {
		var res int64
		err = bus.Object("org.gnome.Mutter.IdleMonitor", "/org/gnome/Mutter/IdleMonitor/Core").Call("org.gnome.Mutter.IdleMonitor.GetIdletime", 0).Store(&res)
		if err != nil {
			return err
		}
		secondsIdle := res / 1000
		if secondsIdle > IDLE_TIME {
			if !wasidle {
				wasidle = true
				onidle()
			}
		} else {
			wasidle = false
		}

		time.Sleep(1 * time.Second)
	}
}
