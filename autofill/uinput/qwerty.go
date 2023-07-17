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
		keyboard.KeyPress(uinput.KeyA)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyB:
		keyboard.KeyPress(uinput.KeyB)
		break
	case KeyBUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyB)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyC:
		keyboard.KeyPress(uinput.KeyC)
		break
	case KeyCUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyC)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyD:
		keyboard.KeyPress(uinput.KeyD)
		break
	case KeyDUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyD)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyE:
		keyboard.KeyPress(uinput.KeyE)
		break
	case KeyEUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyE)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyF:
		keyboard.KeyPress(uinput.KeyF)
		break
	case KeyFUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyF)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyG:
		keyboard.KeyPress(uinput.KeyG)
		break
	case KeyGUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyG)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyH:
		keyboard.KeyPress(uinput.KeyH)
		break
	case KeyHUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyH)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyI:
		keyboard.KeyPress(uinput.KeyI)
		break
	case KeyIUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyI)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyJ:
		keyboard.KeyPress(uinput.KeyJ)
		break
	case KeyJUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyJ)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyK:
		keyboard.KeyPress(uinput.KeyK)
		break
	case KeyKUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyK)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyL:
		keyboard.KeyPress(uinput.KeyL)
		break
	case KeyLUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyL)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyM:
		keyboard.KeyPress(uinput.KeyM)
		break
	case KeyMUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyM)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyN:
		keyboard.KeyPress(uinput.KeyN)
		break
	case KeyNUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyN)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyO:
		keyboard.KeyPress(uinput.KeyO)
		break
	case KeyOUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyO)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyP:
		keyboard.KeyPress(uinput.KeyP)
		break
	case KeyPUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyP)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyQ:
		keyboard.KeyPress(uinput.KeyQ)
		break
	case KeyQUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyQ)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyR:
		keyboard.KeyPress(uinput.KeyR)
		break
	case KeyRUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyR)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyS:
		keyboard.KeyPress(uinput.KeyS)
		break
	case KeySUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyS)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyT:
		keyboard.KeyPress(uinput.KeyT)
		break
	case KeyTUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyT)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyU:
		keyboard.KeyPress(uinput.KeyU)
		break
	case KeyUUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyU)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyV:
		keyboard.KeyPress(uinput.KeyV)
		break
	case KeyVUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyV)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyW:
		keyboard.KeyPress(uinput.KeyW)
		break
	case KeyWUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyW)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyX:
		keyboard.KeyPress(uinput.KeyX)
		break
	case KeyXUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyX)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyY:
		keyboard.KeyPress(uinput.KeyY)
		break
	case KeyYUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyY)
		keyboard.KeyUp(uinput.KeyLeftshift)
	case KeyZ:
		keyboard.KeyPress(uinput.KeyZ)
		break
	case KeyZUpper:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.KeyZ)
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
	case KeyHyphen:
		keyboard.KeyPress(uinput.KeyMinus)
		break
	case KeyTab:
		keyboard.KeyPress(uinput.KeyTab)
		break
	case KeyExclamationMark:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.Key1)
		keyboard.KeyUp(uinput.KeyLeftshift)
		break
	case KeyAtSign:
		keyboard.KeyDown(uinput.KeyLeftshift)
		keyboard.KeyPress(uinput.Key2)
		keyboard.KeyUp(uinput.KeyLeftshift)
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
