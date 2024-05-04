//go:build linux

package notify

import (
	"time"

	"github.com/rymdport/portal/notification"
)

func notifyLibPortal(title string, body string, actionName string, timeout time.Duration, onclose func()) error {
	notificationID++
	err := notification.Add(uint(notificationID), notification.Content{
		Title: title,
		Body:  body,
	})
	if err != nil {
		return err
	}

	if timeout == 0 {
		return nil
	} else {
		go func(id uint32) {
			time.Sleep(timeout)
			notification.Remove(uint(notificationID))
		}(notificationID)
	}
	return nil
}
