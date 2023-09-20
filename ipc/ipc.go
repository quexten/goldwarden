package ipc

import (
	"github.com/quexten/goldwarden/ipc/messages"
)

func ParsedPayload(m messages.IPCMessage) interface{} {
	payload := messages.ParsePayload(m)
	return payload
}

func IPCMessageFromPayload(payload interface{}) (messages.IPCMessage, error) {
	message, err := messages.IPCMessageFromPayload(payload)
	return message, err
}
