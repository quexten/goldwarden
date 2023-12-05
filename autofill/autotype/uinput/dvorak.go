package uinput

import (
	"errors"
	"fmt"

	"github.com/bendahl/uinput"
)

type Dvorak struct {
}

func (d Dvorak) TypeKey(key Key, keyboard uinput.Keyboard) error {
	var err error
	switch key {
	case KeyA:
		err = keyboard.KeyPress(uinput.KeyA)
		break
	case KeyAUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyA)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyB:
		err = keyboard.KeyPress(uinput.KeyN)
		break
	case KeyBUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyN)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyC:
		err = keyboard.KeyPress(uinput.KeyI)
		break
	case KeyCUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyI)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyD:
		err = keyboard.KeyPress(uinput.KeyH)
		break
	case KeyDUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyH)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyE:
		err = keyboard.KeyPress(uinput.KeyD)
		break
	case KeyEUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		err = keyboard.KeyPress(uinput.KeyD)
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyF:
		err = keyboard.KeyPress(uinput.KeyY)
		break
	case KeyFUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyY)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyG:
		err = keyboard.KeyPress(uinput.KeyU)
		break
	case KeyGUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyU)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyH:
		err = keyboard.KeyPress(uinput.KeyJ)
		break
	case KeyHUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyJ)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyI:
		err = keyboard.KeyPress(uinput.KeyG)
		break
	case KeyIUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyG)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyJ:
		err = keyboard.KeyPress(uinput.KeyC)
		break
	case KeyJUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyC)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyK:
		err = keyboard.KeyPress(uinput.KeyV)
		break
	case KeyKUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyV)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyL:
		err = keyboard.KeyPress(uinput.KeyP)
		break
	case KeyLUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyP)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyM:
		err = keyboard.KeyPress(uinput.KeyM)
		break
	case KeyMUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyM)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyN:
		err = keyboard.KeyPress(uinput.KeyL)
		break
	case KeyNUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyL)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyO:
		err = keyboard.KeyPress(uinput.KeyS)
		break
	case KeyOUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyS)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyP:
		err = keyboard.KeyPress(uinput.KeyR)
		break
	case KeyPUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyR)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyQ:
		err = keyboard.KeyPress(uinput.KeyX)
		break
	case KeyQUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyX)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyR:
		err = keyboard.KeyPress(uinput.KeyO)
		break
	case KeyRUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyO)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyS:
		err = keyboard.KeyPress(uinput.KeySemicolon)
		break
	case KeySUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeySemicolon)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyT:
		err = keyboard.KeyPress(uinput.KeyK)
		break
	case KeyTUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyK)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyU:
		err = keyboard.KeyPress(uinput.KeyF)
		break
	case KeyUUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyF)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyV:
		err = keyboard.KeyPress(uinput.KeyDot)
		break
	case KeyVUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyDot)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyW:
		err = keyboard.KeyPress(uinput.KeyComma)
		break
	case KeyWUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyComma)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyX:
		err = keyboard.KeyPress(uinput.KeyB)
		break
	case KeyXUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyB)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyY:
		err = keyboard.KeyPress(uinput.KeyT)
		break
	case KeyYUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyT)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyZ:
		err = keyboard.KeyPress(uinput.KeySlash)
		break
	case KeyZUpper:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeySlash)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
	case Key1:
		err = keyboard.KeyPress(uinput.Key1)
		break
	case Key2:
		err = keyboard.KeyPress(uinput.Key2)
		break
	case Key3:
		err = keyboard.KeyPress(uinput.Key3)
		break
	case Key4:
		err = keyboard.KeyPress(uinput.Key4)
		break
	case Key5:
		err = keyboard.KeyPress(uinput.Key5)
		break
	case Key6:
		err = keyboard.KeyPress(uinput.Key6)
		break
	case Key7:
		err = keyboard.KeyPress(uinput.Key7)
		break
	case Key8:
		err = keyboard.KeyPress(uinput.Key8)
		break
	case Key9:
		err = keyboard.KeyPress(uinput.Key9)
		break
	case Key0:
		err = keyboard.KeyPress(uinput.Key0)
		break
	case KeyTab:
		err = keyboard.KeyPress(uinput.KeyTab)
		break
	case KeyHyphen:
		err = keyboard.KeyPress(uinput.KeyApostrophe)
		break
	case KeyExclamationMark:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key1)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAtSign:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key2)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyHash:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key3)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyDollar:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key4)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyPercent:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key5)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyCaret:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key6)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAmpersand:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key7)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAsterisk:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.Key8)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyDot:
		err = keyboard.KeyPress(uinput.KeyE)
		break
	case KeyComma:
		err = keyboard.KeyPress(uinput.KeyW)
		break
	case KeyQuestionMark:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyLeftbrace)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeySemicolon:
		err = keyboard.KeyPress(uinput.KeyZ)
		break
	case KeyColon:
		err = keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		err = keyboard.KeyPress(uinput.KeyZ)
		Sleep()
		err = keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeySlash:
		err = keyboard.KeyPress(uinput.KeyLeftbrace)
		break
	case KeyApostrophe:
		err = keyboard.KeyPress(uinput.KeyQ)
		break
	case KeySpace:
		err = keyboard.KeyPress(uinput.KeySpace)
		break

	default:
		fmt.Println("Unknown key: ", key)
		fmt.Println("Please add it to the dvorak layout")
		return errors.New("Unknown key")
	}

	return err
}

func init() {
	DefaultLayoutRegistry.Register("dvorak", Dvorak{})
}
