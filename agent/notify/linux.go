//go:build linux || freebsd

package notify

import (
	"time"

	"github.com/quexten/goldwarden/logging"
)

var notificationID uint32 = 1000000
var log = logging.GetLogger("Goldwarden", "Dbus")

func Notify(title string, body string, actionName string, timeout time.Duration, onclose func()) {
	err := notifyLibPortal(title, body, actionName, timeout, onclose)
	if err != nil {
		err = notifyDBus(title, body, actionName, timeout, onclose)
		if err != nil {
			log.Warn("error sending notification " + err.Error())
		}
	}
}

func ListenForNotifications() error {
	return nil
}
