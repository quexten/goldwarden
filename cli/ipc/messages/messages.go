package messages

import (
	"encoding/json"
	"errors"
	"hash/fnv"
	"reflect"

	"github.com/quexten/goldwarden/cli/logging"
)

var log = logging.GetLogger("Goldwarden", "IPC Messages")

type IPCMessageType int64
type IPCMessage struct {
	Type    IPCMessageType `json:"type"`
	Payload []byte         `json:"payload"`
}

type parsePayload func([]byte) (interface{}, error)

var messages = map[string]parsePayload{}
var messageTypes = map[IPCMessageType]string{}

func MessageTypeForEmptyPayload(emptyPayload interface{}) IPCMessageType {
	return hash(reflect.TypeOf(emptyPayload).Name())
}

func hash(s string) IPCMessageType {
	h := fnv.New64()
	h.Write([]byte(s))
	return IPCMessageType(h.Sum64())
}

func registerPayloadParser(payloadParser parsePayload, emptyPayload interface{}) {
	messages[reflect.TypeOf(emptyPayload).Name()] = payloadParser
	messageTypes[hash(reflect.TypeOf(emptyPayload).Name())] = reflect.TypeOf(emptyPayload).Name()
}

func ParsePayload(message IPCMessage) interface{} {
	if _, ok := messageTypes[message.Type]; !ok {
		log.Error("Unregistered message type %d", int(message.Type))
		return nil
	}
	if payload, err := messages[messageTypes[message.Type]](message.Payload); err != nil {
		log.Error("Error parsing payload: %s", err.Error())
		return nil
	} else {
		return payload
	}
}

func IPCMessageFromPayload(payload interface{}) (IPCMessage, error) {
	payloadTypeName := reflect.TypeOf(payload).Name()
	if _, ok := messages[payloadTypeName]; !ok {
		return IPCMessage{}, errors.New("Unregistered payload type " + payloadTypeName)
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return IPCMessage{}, err
	}

	messageType := hash(payloadTypeName)
	return IPCMessage{
		Type:    messageType,
		Payload: payloadJSON,
	}, nil
}
