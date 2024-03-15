package browserbiometrics

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"unsafe"

	"github.com/quexten/goldwarden/browserbiometrics/logging"
)

const bufferSize = 8192 * 8

var nativeEndian binary.ByteOrder

func setupCommunication() {
	// determine native endianess
	var one int16 = 1
	b := (*byte)(unsafe.Pointer(&one))
	if *b == 0 {
		nativeEndian = binary.BigEndian
	} else {
		nativeEndian = binary.LittleEndian
	}
}

func dataToBytes(msg SendMessage) ([]byte, error) {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal OutgoingMessage struct to slice of bytes: %w", err)
	}
	return byteMsg, nil
}

func writeMessageLength(msg []byte) error {
	err := binary.Write(os.Stdout, nativeEndian, uint32(len(msg)))
	if err != nil {
		return fmt.Errorf("unable to write message length to stdout: %w", err)
	}
	return nil
}

func readMessageLength(msg []byte) (int, error) {
	var length int
	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, nativeEndian, &length)
	if err != nil {
		return 0, fmt.Errorf("Unable to read bytes representing message length: %w", err)
	}
	return length, nil
}

func send(msg SendMessage) error {
	byteMsg, err := dataToBytes(msg)
	if err != nil {
		return err
	}

	logging.Debugf("[SENSITIVE] Sending message: " + string(byteMsg))
	err = writeMessageLength(byteMsg)
	if err != nil {
		return err
	}

	var msgBuf bytes.Buffer
	_, err = msgBuf.Write(byteMsg)
	if err != nil {
		return err
	}

	_, err = msgBuf.WriteTo(os.Stdout)
	return err
}
