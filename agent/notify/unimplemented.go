//go:build windows || darwin

package notify

func Notify(title string, body string, actionName string, onclose func()) error {
	// no notifications on windows or darwin
	return nil
}

func ListenForNotifications() error {
	return nil
}
