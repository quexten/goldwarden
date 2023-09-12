//go:build linux

package autotype

import "github.com/quexten/goldwarden/autofill/autotype/uinput"

func TypeString(text string, layout string) error {
	return uinput.TypeString(text, layout)
}

func Paste(layout string) error {
	return uinput.Paste(layout)
}
