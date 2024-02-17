package client

import (
	"encoding/json"
	"io"
	"net"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/ipc/messages"
)

const READ_BUFFER = 4 * 1024 * 1024 // 16MB

type UnixSocketClient struct {
	runtimeConfig *config.RuntimeConfig
}

type UnixSocketConnection struct {
	conn net.Conn
}

func NewUnixSocketClient(runtimeConfig *config.RuntimeConfig) UnixSocketClient {
	return UnixSocketClient{
		runtimeConfig: runtimeConfig,
	}
}

func Reader(r io.Reader) interface{} {
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
	c, err := client.Connect()
	if err != nil {
		return nil, err
	}
	defer c.Close()
	return c.SendCommand(request)
}

func (client UnixSocketClient) Connect() (UnixSocketConnection, error) {
	c, err := net.Dial("unix", client.runtimeConfig.GoldwardenSocketPath)
	if err != nil {
		return UnixSocketConnection{}, err
	}
	return UnixSocketConnection{conn: c}, nil
}

func (conn UnixSocketConnection) SendCommand(request interface{}) (interface{}, error) {
	err := conn.WriteMessage(request)
	if err != nil {
		return nil, err
	}
	return conn.ReadMessage(), nil
}

func (conn UnixSocketConnection) ReadMessage() interface{} {
	result := Reader(conn.conn)
	payload := messages.ParsePayload(result.(messages.IPCMessage))
	return payload
}

func (conn UnixSocketConnection) WriteMessage(message interface{}) error {
	messagePacket, err := messages.IPCMessageFromPayload(message)
	if err != nil {
		panic(err)
	}
	messageJson, err := json.Marshal(messagePacket)
	if err != nil {
		panic(err)
	}
	_, err = conn.conn.Write(messageJson)
	return err
}

func (conn UnixSocketConnection) Close() {
	conn.conn.Close()
}
