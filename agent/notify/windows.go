//go:build windows

package notify

import (
	"time"

	"gopkg.in/toast.v1"
)

func Notify(title string, body string, actionName string, timeout time.Duration, onclose func()) error {
	notification := toast.Notification{
		AppID:   "Goldwarden",
		Title:   title,
		Message: body,
		Audio:   toast.Silent,
	}

	return notification.Push()
}

func ListenForNotifications() error {
	return nil
}
