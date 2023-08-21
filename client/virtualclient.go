package client

import (
	"github.com/quexten/goldwarden/ipc"
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
		message, err := ipc.UnmarshalJSON(<-recv)
		if err != nil {
			panic(err)
		}
		return message
	}
}

func (client VirtualClient) SendToAgent(request interface{}) (interface{}, error) {
	message, err := ipc.IPCMessageFromPayload(request)
	if err != nil {
		panic(err)
	}
	messageJson, err := message.MarshallToJson()
	if err != nil {
		panic(err)
	}

	client.send <- messageJson
	result := virtualReader(client.recv)
	return result.(ipc.IPCMessage).ParsedPayload(), nil
}
