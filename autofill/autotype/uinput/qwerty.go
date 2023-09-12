package uinput

import (
	"errors"
	"fmt"

	"github.com/bendahl/uinput"
)

type Qwerty struct {
}

func (d Qwerty) TypeKey(key Key, keyboard uinput.Keyboard) error {
	switch key {
	case KeyA:
		keyboard.KeyPress(uinput.KeyA)
		break
	case KeyAUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyA)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyB:
		keyboard.KeyPress(uinput.KeyB)
		break
	case KeyBUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyB)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyC:
		keyboard.KeyPress(uinput.KeyC)
		break
	case KeyCUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyC)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyD:
		keyboard.KeyPress(uinput.KeyD)
		break
	case KeyDUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyD)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyE:
		keyboard.KeyPress(uinput.KeyE)
		break
	case KeyEUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyE)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyF:
		keyboard.KeyPress(uinput.KeyF)
		break
	case KeyFUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyF)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyG:
		keyboard.KeyPress(uinput.KeyG)
		break
	case KeyGUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyG)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyH:
		keyboard.KeyPress(uinput.KeyH)
		break
	case KeyHUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyH)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyI:
		keyboard.KeyPress(uinput.KeyI)
		break
	case KeyIUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyI)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyJ:
		keyboard.KeyPress(uinput.KeyJ)
		break
	case KeyJUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyJ)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyK:
		keyboard.KeyPress(uinput.KeyK)
		break
	case KeyKUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyK)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyL:
		keyboard.KeyPress(uinput.KeyL)
		break
	case KeyLUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyL)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyM:
		keyboard.KeyPress(uinput.KeyM)
		break
	case KeyMUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyM)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyN:
		keyboard.KeyPress(uinput.KeyN)
		break
	case KeyNUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyN)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyO:
		keyboard.KeyPress(uinput.KeyO)
		break
	case KeyOUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyO)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyP:
		keyboard.KeyPress(uinput.KeyP)
		break
	case KeyPUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyP)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyQ:
		keyboard.KeyPress(uinput.KeyQ)
		break
	case KeyQUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyQ)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyR:
		keyboard.KeyPress(uinput.KeyR)
		break
	case KeyRUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyR)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyS:
		keyboard.KeyPress(uinput.KeyS)
		break
	case KeySUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyS)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyT:
		keyboard.KeyPress(uinput.KeyT)
		break
	case KeyTUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyT)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyU:
		keyboard.KeyPress(uinput.KeyU)
		break
	case KeyUUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyU)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyV:
		keyboard.KeyPress(uinput.KeyV)
		break
	case KeyVUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyV)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyW:
		keyboard.KeyPress(uinput.KeyW)
		break
	case KeyWUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyW)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyX:
		keyboard.KeyPress(uinput.KeyX)
		break
	case KeyXUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyX)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyY:
		keyboard.KeyPress(uinput.KeyY)
		break
	case KeyYUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyY)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyZ:
		keyboard.KeyPress(uinput.KeyZ)
		break
	case KeyZUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeyZ)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
	case Key1:
		keyboard.KeyPress(uinput.Key1)
		break
	case Key2:
		keyboard.KeyPress(uinput.Key2)
		break
	case Key3:
		keyboard.KeyPress(uinput.Key3)
		break
	case Key4:
		keyboard.KeyPress(uinput.Key4)
		break
	case Key5:
		keyboard.KeyPress(uinput.Key5)
		break
	case Key6:
		keyboard.KeyPress(uinput.Key6)
		break
	case Key7:
		keyboard.KeyPress(uinput.Key7)
		break
	case Key8:
		keyboard.KeyPress(uinput.Key8)
		break
	case Key9:
		keyboard.KeyPress(uinput.Key9)
		break
	case Key0:
		keyboard.KeyPress(uinput.Key0)
		break
	case KeyTab:
		keyboard.KeyPress(uinput.KeyTab)
		break
	case KeyHyphen:
		keyboard.KeyPress(uinput.KeyMinus)
		break
	case KeyExclamationMark:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key1)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAtSign:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key2)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyHash:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key3)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyDollar:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key4)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyPercent:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key5)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyCaret:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key6)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAmpersand:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key7)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAsterisk:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.Key8)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyDot:
		keyboard.KeyPress(uinput.KeyDot)
		break
	case KeyComma:
		keyboard.KeyPress(uinput.KeyW)
		break
	case KeyQuestionMark:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeySlash)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeySemicolon:
		keyboard.KeyPress(uinput.KeySemicolon)
		break
	case KeyColon:
		keyboard.KeyDown(uinput.KeyLeftshift)
		Sleep()
		keyboard.KeyPress(uinput.KeySemicolon)
		Sleep()
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeySlash:
		keyboard.KeyPress(uinput.KeySlash)
		break
	case KeyApostrophe:
		keyboard.KeyPress(uinput.KeyApostrophe)
		break

	case KeySpace:
		keyboard.KeyPress(uinput.KeySpace)
		break

	default:
		fmt.Println("Unknown key: ", key)
		fmt.Println("Please add it to the QWERTY layout")
		return errors.New("Unknown key")
	}
	return nil
}

func init() {
	DefaultLayoutRegistry.Register("qwerty", Qwerty{})
}
