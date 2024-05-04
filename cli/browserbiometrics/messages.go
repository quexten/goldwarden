package browserbiometrics

// top level messages
type GenericRecvMessage struct {
	AppID   string      `json:"appId"`
	Message interface{} `json:"message"`
}

type UnencryptedRecvMessage struct {
	AppID   string         `json:"appId"`
	Message PayloadMessage `json:"message"`
}

type EncryptedRecvMessage struct {
	AppID   string          `json:"appId"`
	Message EncryptedString `json:"message"`
}

type ReceiveMessage struct {
	Timestamp int64  `json:"timestamp"`
	Command   string `json:"command"`
	Response  string `json:"response"`
	KeyB64    string `json:"keyB64"`
}

type SendMessage struct {
	Command      string          `json:"command"`
	AppID        string          `json:"appId"`
	SharedSecret string          `json:"sharedSecret"`
	Message      EncryptedString `json:"message"`
}

type EncryptedString struct {
	IV      string `json:"iv"`
	Mac     string `json:"mac"`
	Data    string `json:"data"`
	EncType int    `json:"encryptionType"`
}

type PayloadMessage struct {
	Command   string `json:"command"`
	UserId    string `json:"userId"`
	Timestamp int64  `json:"timestamp"`
	PublicKey string `json:"publicKey"`
}
