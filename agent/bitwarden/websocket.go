package bitwarden

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/awnumar/memguard"
	"github.com/gorilla/websocket"
	"github.com/quexten/goldwarden/agent/bitwarden/models"
	"github.com/quexten/goldwarden/agent/config"
	"github.com/quexten/goldwarden/agent/systemauth"
	"github.com/quexten/goldwarden/agent/vault"
	"github.com/quexten/goldwarden/logging"
	"github.com/vmihailenco/msgpack/v5"
)

var websocketLog = logging.GetLogger("Goldwarden", "Websocket")

type NotificationMessageType int64

const (
	SyncCipherUpdate NotificationMessageType = 0
	SyncCipherCreate NotificationMessageType = 1
	SyncLoginDelete  NotificationMessageType = 2
	SyncFolderDelete NotificationMessageType = 3
	SyncCiphers      NotificationMessageType = 4

	SyncVault        NotificationMessageType = 5
	SyncOrgKeys      NotificationMessageType = 6
	SyncFolderCreate NotificationMessageType = 7
	SyncFolderUpdate NotificationMessageType = 8
	SyncCipherDelete NotificationMessageType = 9
	SyncSettings     NotificationMessageType = 10

	LogOut NotificationMessageType = 11

	SyncSendCreate NotificationMessageType = 12
	SyncSendUpdate NotificationMessageType = 13
	SyncSendDelete NotificationMessageType = 14

	AuthRequest         NotificationMessageType = 15
	AuthRequestResponse NotificationMessageType = 16
)

const (
	WEBSOCKET_SLEEP_DURATION_SECONDS = 5
)

func RunWebsocketDaemon(ctx context.Context, vault *vault.Vault, cfg *config.Config) {
	for {
		time.Sleep(WEBSOCKET_SLEEP_DURATION_SECONDS * time.Second)

		if cfg.IsLocked() {
			continue
		}

		if token, err := cfg.GetToken(); err == nil && token.AccessToken != "" {
			err := connectToWebsocket(ctx, vault, cfg)
			if err != nil {
				websocketLog.Error("Websocket error %s", err)
			}
		}
	}
}

func connectToWebsocket(ctx context.Context, vault *vault.Vault, cfg *config.Config) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	url, err := url.Parse(cfg.ConfigFile.NotificationsUrl)
	if err != nil {
		return err
	}

	token, err := cfg.GetToken()
	var websocketURL = "wss://" + url.Host + url.Path + "/hub?access_token=" + token.AccessToken
	c, _, err := websocket.DefaultDialer.Dial(websocketURL, nil)
	if err != nil {
		return err
	}
	defer c.Close()

	websocketLog.Info("Connected to websocket server...")

	done := make(chan struct{})
	//handshake required for official bitwarden implementation
	c.WriteMessage(1, []byte(`{"protocol":"messagepack","version":1}`))

	go func() {
		defer close(done)
		for {
			mt, message, err := c.ReadMessage()
			fmt.Println(mt)
			if err != nil {
				websocketLog.Error("Error reading websocket message %s", err)
				return
			}

			if messageType, cipherid, success := websocketMessageType(message); success {
				var mt1 = NotificationMessageType(messageType)
				switch mt1 {
				case SyncCiphers, SyncVault:
					websocketLog.Warn("SyncCiphers requested")
					token, err := cfg.GetToken()
					if err != nil {
						websocketLog.Error("Error getting token %s", err)
						break
					}
					DoFullSync(context.WithValue(ctx, AuthToken{}, token.AccessToken), vault, cfg, nil, false)
					break
				case SyncCipherDelete:
					websocketLog.Warn("Delete requested for cipher " + cipherid)
					vault.DeleteCipher(cipherid)
					break
				case SyncCipherUpdate:
					websocketLog.Warn("Update requested for cipher " + cipherid)
					token, err := cfg.GetToken()
					if err != nil {
						websocketLog.Error("Error getting token %s", err)
						break
					}

					cipher, err := GetCipher(context.WithValue(ctx, AuthToken{}, token.AccessToken), cipherid, cfg)
					if err != nil {
						websocketLog.Error("Error getting cipher %s", err)
						break
					}
					if !cipher.DeletedDate.IsZero() {
						websocketLog.Info("Cipher moved to trash " + cipherid)
						vault.DeleteCipher(cipherid)
						break
					}

					if cipher.Type == models.CipherNote {
						vault.AddOrUpdateSecureNote(cipher)
					} else {
						vault.AddOrUpdateLogin(cipher)
					}
					break
				case SyncCipherCreate:
					websocketLog.Warn("Create requested for cipher " + cipherid)
					token, err := cfg.GetToken()
					if err != nil {
						websocketLog.Error("Error getting token %s", err)
						break
					}

					cipher, err := GetCipher(context.WithValue(ctx, AuthToken{}, token.AccessToken), cipherid, cfg)
					if err != nil {
						websocketLog.Error("Error getting cipher %s", err)
						break
					}

					if cipher.Type == models.CipherNote {
						vault.AddOrUpdateSecureNote(cipher)
					} else {
						vault.AddOrUpdateLogin(cipher)
					}

					break
				case SyncSendCreate, SyncSendUpdate, SyncSendDelete:
					websocketLog.Warn("SyncSend requested: sends are not supported")
					break
				case LogOut:
					websocketLog.Info("LogOut received. Wiping vault and exiting...")
					memguard.SafeExit(0)
				case AuthRequest:
					websocketLog.Info("AuthRequest received" + string(cipherid))
					authRequest, err := GetAuthRequest(context.WithValue(ctx, AuthToken{}, token.AccessToken), cipherid, cfg)
					if err != nil {
						websocketLog.Error("Error getting auth request %s", err)
						break
					}
					websocketLog.Info("AuthRequest details " + authRequest.RequestIpAddress + " " + authRequest.RequestDeviceType)

					if approved, err := systemauth.GetApproval("Paswordless Login Request", "Do you want to allow "+authRequest.RequestIpAddress+" ("+authRequest.RequestDeviceType+") to login to your account?"); err != nil || !approved {
						websocketLog.Info("AuthRequest denied")
						break
					}
					if !systemauth.CheckBiometrics(systemauth.AccessCredential) {
						websocketLog.Info("AuthRequest denied - biometrics required")
						break
					}

					_, err = CreateAuthResponse(context.WithValue(ctx, AuthToken{}, token.AccessToken), authRequest, vault.Keyring, cfg)
					if err != nil {
						websocketLog.Error("Error creating auth response %s", err)
					}
					break
				case AuthRequestResponse:
					websocketLog.Info("AuthRequestResponse received")
					break
				case SyncFolderDelete, SyncFolderCreate, SyncFolderUpdate:
					websocketLog.Warn("SyncFolder requested: folders are not supported")
					break
				case SyncOrgKeys, SyncSettings:
					websocketLog.Warn("SyncOrgKeys requested: orgs / settings are not supported")
					break
				default:
					websocketLog.Warn("Unknown message type received %d", mt1)
				}
			}
		}
	}()

	<-done
	return nil
}

func websocketMessageType(message []byte) (int8, string, bool) {
	lenBufferLen := 0
	for i := 0; i < len(message); i++ {
		if (message[i] & 0x80) == 0 {
			lenBufferLen = i + 1
			break
		}
	}
	msgPackMessage := message[lenBufferLen:]
	return parseMessageTypeFromMessagePack(msgPackMessage)
}

func parseMessageTypeFromMessagePack(messagePack []byte) (int8, string, bool) {
	msgPackBuffer := bytes.NewBuffer(messagePack)
	dec := msgpack.NewDecoder(msgPackBuffer)
	value, err := dec.DecodeSlice()
	if value == nil || err != nil {
		return 0, "", false
	}
	if len(value) < 5 {
		websocketLog.Warn("Invalid message received, length too short")
		return 0, "", false
	}
	value, success := value[4].([]interface{})
	if len(value) < 1 || !success {
		websocketLog.Warn("Invalid message received, length too short")
		return 0, "", false
	}
	value1, success := value[0].(map[string]interface{})
	if !success {
		websocketLog.Warn("Invalid message received, value is not a map")
		return 0, "", false
	}
	if _, ok := value1["Type"]; !ok {
		websocketLog.Warn("Invalid message received, no type")
		return 0, "", false
	}
	messagePayloadType, success := value1["Type"].(int8)
	if !success {
		websocketLog.Warn("Invalid message received, type is not an int")
		return 0, "", false
	}
	payload, success := value1["Payload"].(map[string]interface{})
	if !success {
		return messagePayloadType, "", true
	}
	if _, ok := payload["Id"]; !ok {
		websocketLog.Warn("Invalid message received, no id")
		return 0, "", false
	}

	return messagePayloadType, payload["Id"].(string), true
}
