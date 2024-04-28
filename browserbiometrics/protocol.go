package browserbiometrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/browserbiometrics/logging"
	"github.com/quexten/goldwarden/client"
	"github.com/quexten/goldwarden/ipc/messages"
)

var runtimeConfig *config.RuntimeConfig

func readLoop(rtCfg *config.RuntimeConfig) error {
	runtimeConfig = rtCfg
	v := bufio.NewReader(os.Stdin)
	s := bufio.NewReaderSize(v, bufferSize)

	lengthBytes := make([]byte, 4)
	lengthNum := int(0)

	logging.Debugf("Sending connected message")
	err := send(SendMessage{
		Command: "connected",
		AppID:   appID,
	})
	if err != nil {
		return err
	}

	logging.Debugf("Starting read loop")
	for b, err := s.Read(lengthBytes); b > 0 && err == nil; b, err = s.Read(lengthBytes) {
		lengthNum, err = readMessageLength(lengthBytes)
		if err != nil {
			return err
		}

		content := make([]byte, lengthNum)
		_, err := s.Read(content)
		if err != nil && err != io.EOF {
			return err
		}

		err = parseMessage(content)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseMessage(msg []byte) error {
	logging.Debugf("Received message: " + string(msg))

	var genericMessage GenericRecvMessage
	err := json.Unmarshal(msg, &genericMessage)
	if err != nil {
		return fmt.Errorf("unable to unmarshal json to struct: %w", err)
	}

	if _, ok := (genericMessage.Message.(map[string]interface{})["command"]); ok {
		logging.Debugf("Message is unencrypted")

		var unmsg UnencryptedRecvMessage
		err := json.Unmarshal(msg, &unmsg)
		if err != nil {
			return fmt.Errorf("unable to unmarshal json to struct: %w", err)
		}

		err = handleUnencryptedMessage(unmsg)
		if err != nil {
			return err
		}
	} else {
		logging.Debugf("Message is encrypted")

		var encmsg EncryptedRecvMessage
		err := json.Unmarshal(msg, &encmsg)
		if err != nil {
			return fmt.Errorf("unable to unmarshal json to struct: %w", err)
		}

		decryptedMessage, err := decryptStringSymmetric(transportKey, encmsg.Message.IV, encmsg.Message.Data)
		if err != nil {
			return err
		}
		var payloadMsg PayloadMessage
		err = json.Unmarshal([]byte(decryptedMessage), &payloadMsg)
		if err != nil {
			return fmt.Errorf("unable to unmarshal json to struct: %w", err)
		}

		err = handlePayloadMessage(payloadMsg, genericMessage.AppID)
		if err != nil {
			return err
		}
	}

	return nil
}

func handleUnencryptedMessage(msg UnencryptedRecvMessage) error {
	logging.Debugf("Received unencrypted message: %+v", msg.Message)
	logging.Debugf("  with command: %s", msg.Message.Command)

	switch msg.Message.Command {
	case "setupEncryption":
		sharedSecret, err := rsaEncrypt(msg.Message.PublicKey, transportKey)
		if err != nil {
			return err
		}
		err = send(SendMessage{
			Command:      "setupEncryption",
			AppID:        msg.AppID,
			SharedSecret: sharedSecret,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func handlePayloadMessage(msg PayloadMessage, appID string) error {
	logging.Debugf("Received unencrypted message: %+v", msg)

	switch msg.Command {
	case "biometricUnlock":
		logging.Debugf("Biometric unlock requested")
		// logging.Debugf("Biometrics authorized: %t", isAuthorized)

		home, err := os.UserHomeDir()
		if err != nil {
			return err
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
			return fmt.Errorf("Unable to send message to agent: %w", err)
		}

		switch result := result.(type) {
		case messages.GetBiometricsKeyResponse:
			var key = result.Key
			var payloadMsg ReceiveMessage = ReceiveMessage{
				Command:   "biometricUnlock",
				Response:  "unlocked",
				Timestamp: msg.Timestamp,
				KeyB64:    key,
			}
			payloadStr, err := json.Marshal(payloadMsg)
			if err != nil {
				return err
			}
			logging.Debugf("Payload: %s", payloadStr)

			encStr, err := encryptStringSymmetric(transportKey, payloadStr)
			if err != nil {
				return err
			}
			err = send(SendMessage{
				AppID:   appID,
				Message: encStr,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
