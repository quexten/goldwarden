//go:build darwin

package notify

import (
	"time"

	"github.com/gen2brain/beeep"
)

func Notify(title string, body string, actionName string, timeout time.Duration, onclose func()) error {
	err := beeep.Notify(title, body, "")
	if err != nil {
		panic(err)
	}
	return nil
}

func ListenForNotifications() error {
	return nil
}
