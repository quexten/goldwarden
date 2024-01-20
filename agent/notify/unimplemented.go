//go:build windows || darwin

package notify

import "time"

func Notify(title string, body string, actionName string, timeout time.Duration, onclose func()) error {
	// no notifications on windows or darwin
	return nil
}

func ListenForNotifications() error {
	return nil
}
