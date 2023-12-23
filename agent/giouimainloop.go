//go:build (windows || darwin || linux) && !noautofill

package agent

import "gioui.org/app"

func mainloop() {
	app.Main()
}
