package uinput

import (
	"errors"
	"fmt"
	"time"

	"github.com/bendahl/uinput"
)

type Layout interface {
	TypeKey(key Key, keyboard uinput.Keyboard) error
}

type Key string

const (
	KeyA      Key = "a"
	KeyB      Key = "b"
	KeyC      Key = "c"
	KeyD      Key = "d"
	KeyE      Key = "e"
	KeyF      Key = "f"
	KeyG      Key = "g"
	KeyH      Key = "h"
	KeyI      Key = "i"
	KeyJ      Key = "j"
	KeyK      Key = "k"
	KeyL      Key = "l"
	KeyM      Key = "m"
	KeyN      Key = "n"
	KeyO      Key = "o"
	KeyP      Key = "p"
	KeyQ      Key = "q"
	KeyR      Key = "r"
	KeyS      Key = "s"
	KeyT      Key = "t"
	KeyU      Key = "u"
	KeyV      Key = "v"
	KeyW      Key = "w"
	KeyX      Key = "x"
	KeyY      Key = "y"
	KeyZ      Key = "z"
	KeyAUpper Key = "A"
	KeyBUpper Key = "B"
	KeyCUpper Key = "C"
	KeyDUpper Key = "D"
	KeyEUpper Key = "E"
	KeyFUpper Key = "F"
	KeyGUpper Key = "G"
	KeyHUpper Key = "H"
	KeyIUpper Key = "I"
	KeyJUpper Key = "J"
	KeyKUpper Key = "K"
	KeyLUpper Key = "L"
	KeyMUpper Key = "M"
	KeyNUpper Key = "N"
	KeyOUpper Key = "O"
	KeyPUpper Key = "P"
	KeyQUpper Key = "Q"
	KeyRUpper Key = "R"
	KeySUpper Key = "S"
	KeyTUpper Key = "T"
	KeyUUpper Key = "U"
	KeyVUpper Key = "V"
	KeyWUpper Key = "W"
	KeyXUpper Key = "X"
	KeyYUpper Key = "Y"
	KeyZUpper Key = "Z"
	Key0      Key = "0"
	Key1      Key = "1"
	Key2      Key = "2"
	Key3      Key = "3"
	Key4      Key = "4"
	Key5      Key = "5"
	Key6      Key = "6"
	Key7      Key = "7"
	Key8      Key = "8"
	Key9      Key = "9"

	KeyHyphen Key = "-"

	KeySpace           Key = " "
	KeyExclamationMark Key = "!"
	KeyAtSign          Key = "@"
	KeyHash            Key = "#"
	KeyDollar          Key = "$"
	KeyPercent         Key = "%"
	KeyCaret           Key = "^"
	KeyAmpersand       Key = "&"
	KeyAsterisk        Key = "*"

	KeyDot          Key = "."
	KeyComma        Key = ","
	KeySlash        Key = "/"
	KeyBackslash    Key = "\\"
	KeyQuestionMark Key = "?"
	KeySemicolon    Key = ";"
	KeyColon        Key = ":"
	KeyApostrophe   Key = "'"

	KeyTab Key = "\t"
)

type LayoutRegistry struct {
	layouts map[string]Layout
}

func NewLayoutRegistry() *LayoutRegistry {
	return &LayoutRegistry{
		layouts: make(map[string]Layout),
	}
}

var DefaultLayoutRegistry = NewLayoutRegistry()

func (r *LayoutRegistry) Register(name string, layout Layout) {
	r.layouts[name] = layout
}

func TypeString(text string, layout string) error {
	if layout == "" {
		layout = "qwerty"
	}

	if _, ok := DefaultLayoutRegistry.layouts[layout]; !ok {
		return errors.New("layout not found")
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("testkeyboard"))
	if err != nil {
		return err
	}

	for _, c := range text {
		key := Key(string(c))
		err := DefaultLayoutRegistry.layouts[layout].TypeKey(key, keyboard)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(10 * time.Millisecond)
	}

	err = keyboard.Close()
	if err != nil {
		return err
	}

	return nil
}

func Paste(layout string) error {
	if layout == "" {
		layout = "qwerty"
	}

	if _, ok := DefaultLayoutRegistry.layouts[layout]; !ok {
		return errors.New("layout not found")
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("Goldwarden Autotype"))
	if err != nil {
		return err
	}
	keyboard.KeyDown(uinput.KeyLeftctrl)
	time.Sleep(100 * time.Millisecond)
	DefaultLayoutRegistry.layouts[layout].TypeKey(KeyV, keyboard)
	time.Sleep(100 * time.Millisecond)
	keyboard.KeyUp(uinput.KeyLeftctrl)
	return nil
}

func Sleep() {
	time.Sleep(20 * time.Millisecond)
}
