package browserbiometrics

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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

func dataToBytes(msg SendMessage) []byte {
	byteMsg, err := json.Marshal(msg)
	if err != nil {
		logging.Panicf("Unable to marshal OutgoingMessage struct to slice of bytes: " + err.Error())
	}
	return byteMsg
}

func writeMessageLength(msg []byte) {
	err := binary.Write(os.Stdout, nativeEndian, uint32(len(msg)))
	if err != nil {
		logging.Panicf("Unable to write message length to Stdout: " + err.Error())
	}
}

func readMessageLength(msg []byte) int {
	var length uint32
	buf := bytes.NewBuffer(msg)
	err := binary.Read(buf, nativeEndian, &length)
	if err != nil {
		logging.Panicf("Unable to read bytes representing message length:" + err.Error())
	}
	return int(length)
}

func send(msg SendMessage) {
	byteMsg := dataToBytes(msg)
	logging.Debugf("Sending message: " + string(byteMsg))
	writeMessageLength(byteMsg)

	var msgBuf bytes.Buffer
	_, err := msgBuf.Write(byteMsg)
	if err != nil {
		logging.Panicf(err.Error())
	}

	_, err = msgBuf.WriteTo(os.Stdout)
	if err != nil {
		logging.Panicf(err.Error())
	}
}
