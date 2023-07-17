package client

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/quexten/goldwarden/ipc"
)

const READ_BUFFER = 1 * 1024 * 1024 // 1MB

func reader(r io.Reader) interface{} {
	buf := make([]byte, READ_BUFFER)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return nil
		}
		message, err := ipc.UnmarshalJSON(buf[0:n])
		if err != nil {
			panic(err)
		}
		return message
	}
}

func SendToAgent(request interface{}) (interface{}, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	// home := "/home/quexten"
	c, err := net.Dial("unix", home+"/.goldwarden.sock")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	message, err := ipc.IPCMessageFromPayload(request)
	if err != nil {
		panic(err)
	}
	messageJson, err := message.MarshallToJson()
	if err != nil {
		panic(err)
	}

	_, err = c.Write(messageJson)
	if err != nil {
		log.Fatal("write error:", err)
	}
	result := reader(c)
	return result.(ipc.IPCMessage).ParsedPayload(), nil
}
