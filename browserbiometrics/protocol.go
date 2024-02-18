package browserbiometrics

import (
	"bufio"
	"encoding/json"
	"io"
	"os"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/browserbiometrics/logging"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc/messages"
)

var runtimeConfig *config.RuntimeConfig

func readLoop(rtCfg *config.RuntimeConfig) {
	runtimeConfig = rtCfg
	v := bufio.NewReader(os.Stdin)
	s := bufio.NewReaderSize(v, bufferSize)

	lengthBytes := make([]byte, 4)
	lengthNum := int(0)

	logging.Debugf("Sending connected message")
	send(SendMessage{
		Command: "connected",
		AppID:   appID,
	})

	logging.Debugf("Starting read loop")
	for b, err := s.Read(lengthBytes); b > 0 && err == nil; b, err = s.Read(lengthBytes) {
		lengthNum = readMessageLength(lengthBytes)

		content := make([]byte, lengthNum)
		_, err := s.Read(content)
		if err != nil && err != io.EOF {
			logging.Panicf(err.Error())
		}

		parseMessage(content)
	}
}

func parseMessage(msg []byte) {
	logging.Debugf("Received message: " + string(msg))

	var genericMessage GenericRecvMessage
	err := json.Unmarshal(msg, &genericMessage)
	if err != nil {
		logging.Panicf("Unable to unmarshal json to struct: " + err.Error())
	}
	if _, ok := (genericMessage.Message.(map[string]interface{})["command"]); ok {
		logging.Debugf("Message is unencrypted")

		var unmsg UnencryptedRecvMessage
		err := json.Unmarshal(msg, &unmsg)
		if err != nil {
			logging.Panicf("Unable to unmarshal json to struct: " + err.Error())
		}

		handleUnencryptedMessage(unmsg)
	} else {
		logging.Debugf("Message is encrypted")

		var encmsg EncryptedRecvMessage
		err := json.Unmarshal(msg, &encmsg)
		if err != nil {
			logging.Panicf("Unable to unmarshal json to struct: " + err.Error())
		}

		decryptedMessage := decryptStringSymmetric(transportKey, encmsg.Message.IV, encmsg.Message.Data)
		var payloadMsg PayloadMessage
		err = json.Unmarshal([]byte(decryptedMessage), &payloadMsg)
		if err != nil {
			logging.Panicf("Unable to unmarshal json to struct: " + err.Error())
		}

		handlePayloadMessage(payloadMsg, genericMessage.AppID)
	}
}

func handleUnencryptedMessage(msg UnencryptedRecvMessage) {
	logging.Debugf("Received unencrypted message: %+v", msg.Message)
	logging.Debugf("  with command: %s", msg.Message.Command)

	switch msg.Message.Command {
	case "setupEncryption":
		sharedSecret, err := rsaEncrypt(msg.Message.PublicKey, transportKey)
		if err != nil {
			logging.Panicf(err.Error())
		}
		send(SendMessage{
			Command:      "setupEncryption",
			AppID:        msg.AppID,
			SharedSecret: sharedSecret,
		})
		break
	}
}
func handlePayloadMessage(msg PayloadMessage, appID string) {
	logging.Debugf("Received unencrypted message: %+v", msg)

	switch msg.Command {
	case "biometricUnlock":
		logging.Debugf("Biometric unlock requested")
		// logging.Debugf("Biometrics authorized: %t", isAuthorized)

		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		if runtimeConfig.GoldwardenSocketPath == "" {
			if _, err := os.Stat(home + "/.goldwarden.sock"); err == nil {
				runtimeConfig.GoldwardenSocketPath = home + "/.goldwarden.sock"
			} else if _, err := os.Stat(home + "/.var/app/com.quexten.Goldwarden/data/goldwarden.sock"); err == nil {
				runtimeConfig.GoldwardenSocketPath = home + "/.var/app/com.quexten.Goldwarden/data/goldwarden.sock"
			}

			if _, err = os.Stat("/.flatpak-info"); err == nil {
				runtimeConfig.GoldwardenSocketPath = home + "/.var/app/com.quexten.Goldwarden/data/goldwarden.sock"
			}
		}

		logging.Debugf("Connecting to agent at path %s", runtimeConfig.GoldwardenSocketPath)

		result, err := client.NewUnixSocketClient(runtimeConfig).SendToAgent(messages.GetBiometricsKeyRequest{})
		if err != nil {
			logging.Errorf("Unable to send message to agent: %s", err.Error())
			return
		}
		switch result.(type) {
		case messages.GetBiometricsKeyResponse:
			if err != nil {
				logging.Panicf(err.Error())
			}

			var key = result.(messages.GetBiometricsKeyResponse).Key
			var payloadMsg ReceiveMessage = ReceiveMessage{
				Command:   "biometricUnlock",
				Response:  "unlocked",
				Timestamp: msg.Timestamp,
				KeyB64:    key,
			}
			payloadStr, err := json.Marshal(payloadMsg)
			if err != nil {
				logging.Panicf(err.Error())
			}
			logging.Debugf("Payload: %s", payloadStr)

			encStr := encryptStringSymmetric(transportKey, payloadStr)
			send(SendMessage{
				AppID:   appID,
				Message: encStr,
			})
			break
		}

		break
	}
}
