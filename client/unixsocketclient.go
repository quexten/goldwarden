package client

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"os"

	"github.com/quexten/goldwarden/ipc/messages"
)

const READ_BUFFER = 1 * 1024 * 1024 // 1MB

type UnixSocketClient struct {
}

func NewUnixSocketClient() UnixSocketClient {
	return UnixSocketClient{}
}

func reader(r io.Reader) interface{} {
	buf := make([]byte, READ_BUFFER)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return nil
		}

		var message messages.IPCMessage
		err = json.Unmarshal(buf[0:n], &message)
		if err != nil {
			panic(err)
		}
		return message
	}
}

func (client UnixSocketClient) SendToAgent(request interface{}) (interface{}, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	c, err := net.Dial("unix", home+"/.goldwarden.sock")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	message, err := messages.IPCMessageFromPayload(request)
	if err != nil {
		panic(err)
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	_, err = c.Write(messageJson)
	if err != nil {
		log.Fatal("write error:", err)
	}
	result := reader(c)
	return messages.ParsePayload(result.(messages.IPCMessage)), nil
}
