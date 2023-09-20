//go:build windows || darwin || linux

package agent

import "gioui.org/app"

func mainloop() {
	app.Main()
}
