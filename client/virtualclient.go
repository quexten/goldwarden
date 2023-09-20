package client

import (
	"encoding/json"

	"github.com/quexten/goldwarden/ipc/messages"
)

func NewVirtualClient(recv chan []byte, send chan []byte) VirtualClient {
	return VirtualClient{
		recv,
		send,
	}
}

type VirtualClient struct {
	recv chan []byte
	send chan []byte
}

func virtualReader(recv chan []byte) interface{} {
	for {
		var message messages.IPCMessage
		err := json.Unmarshal(<-recv, &message)
		if err != nil {
			panic(err)
		}
		return message
	}
}

func (client VirtualClient) SendToAgent(request interface{}) (interface{}, error) {
	message, err := messages.IPCMessageFromPayload(request)
	if err != nil {
		panic(err)
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	client.send <- messageJson
	result := virtualReader(client.recv)
	return messages.ParsePayload(result.(messages.IPCMessage)), nil
}
