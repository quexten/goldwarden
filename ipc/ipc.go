package ipc

import (
	"encoding/json"
)

type IPCMessageType int64

const (
	IPCMessageTypeErrorMessage IPCMessageType = 0

	IPCMessageTypeDoLoginRequest IPCMessageType = 1

	IPCMessageTypeUpdateVaultPINRequest IPCMessageType = 4
	IPCMessageTypeUnlockVaultRequest    IPCMessageType = 5
	IPCMessageTypeLockVaultRequest      IPCMessageType = 6
	IPCMessageTypeWipeVaultRequest      IPCMessageType = 7

	IPCMessageTypeGetCLICredentialsRequest  IPCMessageType = 11
	IPCMessageTypeGetCLICredentialsResponse IPCMessageType = 12

	IPCMessageTypeCreateSSHKeyRequest  IPCMessageType = 14
	IPCMessageTypeCreateSSHKeyResponse IPCMessageType = 15

	IPCMessageTypeGetSSHKeysRequest  IPCMessageType = 16
	IPCMessageTypeGetSSHKeysResponse IPCMessageType = 17

	IPCMessageGetLoginRequest  IPCMessageType = 18
	IPCMessageGetLoginResponse IPCMessageType = 19

	IPCMessageAddLoginRequest  IPCMessageType = 20
	IPCMessageAddLoginResponse IPCMessageType = 21

	IPCMessageGetNoteRequest    IPCMessageType = 26
	IPCMessageGetNoteResponse   IPCMessageType = 27
	IPCMessageGetNotesResponse  IPCMessageType = 32
	IPCMessageGetLoginsResponse IPCMessageType = 33

	IPCMessageAddNoteRequest  IPCMessageType = 28
	IPCMessageAddNoteResponse IPCMessageType = 29

	IPCMessageListLoginsRequest IPCMessageType = 22

	IPCMessageTypeActionResponse IPCMessageType = 13

	IPCMessageTypeGetVaultPINStatusRequest IPCMessageType = 2

	IPCMessageTypeSetAPIUrlRequest      IPCMessageType = 30
	IPCMessageTypeSetIdentityURLRequest IPCMessageType = 31
)

type IPCMessage struct {
	Type    IPCMessageType `json:"type"`
	Payload []byte         `json:"payload"`
}

func (m IPCMessage) MarshallToJson() ([]byte, error) {
	return json.Marshal(m)
}

func UnmarshalJSON(data []byte) (IPCMessage, error) {
	var m IPCMessage
	err := json.Unmarshal(data, &m)
	return m, err
}

func (m IPCMessage) ParsedPayload() interface{} {
	switch m.Type {
	case IPCMessageTypeDoLoginRequest:
		var req DoLoginRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeActionResponse:
		var res ActionResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageTypeErrorMessage:
		return nil
	case IPCMessageTypeGetCLICredentialsRequest:
		var req GetCLICredentialsRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeGetCLICredentialsResponse:
		var res GetCLICredentialsResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageTypeCreateSSHKeyRequest:
		var req CreateSSHKeyRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeCreateSSHKeyResponse:
		var res CreateSSHKeyResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageTypeGetSSHKeysRequest:
		var req GetSSHKeysRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeGetSSHKeysResponse:
		var res GetSSHKeysResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageGetLoginRequest:
		var req GetLoginRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageGetLoginResponse:
		var res GetLoginResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageAddLoginRequest:
		var req AddLoginRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageAddLoginResponse:
		var res AddLoginResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	case IPCMessageTypeWipeVaultRequest:
		var req WipeVaultRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeLockVaultRequest:
		var req LockVaultRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeGetVaultPINStatusRequest:
		var req GetVaultPINRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeSetAPIUrlRequest:
		var req SetApiURLRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageTypeSetIdentityURLRequest:
		var req SetIdentityURLRequest
		err := json.Unmarshal(m.Payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req
	case IPCMessageGetLoginsResponse:
		var res GetLoginsResponse
		err := json.Unmarshal(m.Payload, &res)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return res
	default:
		return nil
	}
}

func IPCMessageFromPayload(payload interface{}) (IPCMessage, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return IPCMessage{}, err
	}

	switch payload.(type) {
	case UnlockVaultRequest:
		return IPCMessage{
			Type:    IPCMessageTypeUnlockVaultRequest,
			Payload: jsonBytes,
		}, nil
	case UpdateVaultPINRequest:
		return IPCMessage{
			Type:    IPCMessageTypeUpdateVaultPINRequest,
			Payload: jsonBytes,
		}, nil
	case DoLoginRequest:
		return IPCMessage{
			Type:    IPCMessageTypeDoLoginRequest,
			Payload: jsonBytes,
		}, nil
	case ActionResponse:
		return IPCMessage{
			Type:    IPCMessageTypeActionResponse,
			Payload: jsonBytes,
		}, nil
	case GetCLICredentialsRequest:
		return IPCMessage{
			Type:    IPCMessageTypeGetCLICredentialsRequest,
			Payload: jsonBytes,
		}, nil
	case GetCLICredentialsResponse:
		return IPCMessage{
			Type:    IPCMessageTypeGetCLICredentialsResponse,
			Payload: jsonBytes,
		}, nil
	case CreateSSHKeyRequest:
		return IPCMessage{
			Type:    IPCMessageTypeCreateSSHKeyRequest,
			Payload: jsonBytes,
		}, nil
	case CreateSSHKeyResponse:
		return IPCMessage{
			Type:    IPCMessageTypeCreateSSHKeyResponse,
			Payload: jsonBytes,
		}, nil
	case GetSSHKeysRequest:
		return IPCMessage{
			Type:    IPCMessageTypeGetSSHKeysRequest,
			Payload: jsonBytes,
		}, nil
	case GetSSHKeysResponse:
		return IPCMessage{
			Type:    IPCMessageTypeGetSSHKeysResponse,
			Payload: jsonBytes,
		}, nil
	case WipeVaultRequest:
		return IPCMessage{
			Type:    IPCMessageTypeWipeVaultRequest,
			Payload: jsonBytes,
		}, nil
	case LockVaultRequest:
		return IPCMessage{
			Type:    IPCMessageTypeLockVaultRequest,
			Payload: jsonBytes,
		}, nil
	case GetVaultPINRequest:
		return IPCMessage{
			Type:    IPCMessageTypeGetVaultPINStatusRequest,
			Payload: jsonBytes,
		}, nil
	case SetApiURLRequest:
		return IPCMessage{
			Type:    IPCMessageTypeSetAPIUrlRequest,
			Payload: jsonBytes,
		}, nil
	case SetIdentityURLRequest:
		return IPCMessage{
			Type:    IPCMessageTypeSetIdentityURLRequest,
			Payload: jsonBytes,
		}, nil
	case GetLoginRequest:
		return IPCMessage{
			Type:    IPCMessageGetLoginRequest,
			Payload: jsonBytes,
		}, nil
	case GetLoginResponse:
		return IPCMessage{
			Type:    IPCMessageGetLoginResponse,
			Payload: jsonBytes,
		}, nil
	case AddLoginRequest:
		return IPCMessage{
			Type:    IPCMessageAddLoginRequest,
			Payload: jsonBytes,
		}, nil
	case AddLoginResponse:
		return IPCMessage{
			Type:    IPCMessageAddLoginResponse,
			Payload: jsonBytes,
		}, nil
	case GetNotesRequest:
		return IPCMessage{
			Type:    IPCMessageGetNoteRequest,
			Payload: jsonBytes,
		}, nil
	case GetNotesResponse:
		return IPCMessage{
			Type:    IPCMessageGetNotesResponse,
			Payload: jsonBytes,
		}, nil
	case GetNoteResponse:
		return IPCMessage{
			Type:    IPCMessageGetNoteResponse,
			Payload: jsonBytes,
		}, nil
	case GetLoginsResponse:
		return IPCMessage{
			Type:    IPCMessageGetLoginsResponse,
			Payload: jsonBytes,
		}, nil
	case ListLoginsRequest:
		return IPCMessage{
			Type:    IPCMessageListLoginsRequest,
			Payload: jsonBytes,
		}, nil
	default:
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			return IPCMessage{}, err
		}

		return IPCMessage{
			Type:    IPCMessageTypeErrorMessage,
			Payload: payloadBytes,
		}, nil
	}
}

type DoLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LockVaultRequest struct {
}

type UnlockVaultRequest struct {
}

type UpdateVaultPINRequest struct {
}

type ActionResponse struct {
	Success bool
	Message string
}

type GetCLICredentialsRequest struct {
	ApplicationName string
}

type GetCLICredentialsResponse struct {
	Env map[string]string
}

type CreateSSHKeyRequest struct {
	Name string
}

type CreateSSHKeyResponse struct {
	Digest string
}

type GetSSHKeysRequest struct {
}

type GetSSHKeysResponse struct {
	Keys []string
}

type GetLoginRequest struct {
	Name     string
	Username string
	UUID     string
	OrgId    string

	GetList bool
}

type GetLoginResponse struct {
	Found  bool
	Result DecryptedLoginCipher
}

type GetLoginsResponse struct {
	Found  bool
	Result []DecryptedLoginCipher
}

type DecryptedLoginCipher struct {
	Name          string
	Username      string
	Password      string
	UUID          string
	OrgaizationID string
	Notes         string
}

type GetNotesRequest struct {
	Name string
}

type GetNoteResponse struct {
	Found  bool
	Result DecryptedNoteCipher
}

type GetNotesResponse struct {
	Found  bool
	Result []DecryptedNoteCipher
}

type DecryptedNoteCipher struct {
	Name     string
	Contents string
}

type AddLoginRequest struct {
	Name string
	UUID string
}

type AddLoginResponse struct {
	Name string
	UUID string
}

type WipeVaultRequest struct {
}

type GetVaultPINRequest struct {
}

type SetApiURLRequest struct {
	Value string
}

type SetIdentityURLRequest struct {
	Value string
}

type ListLoginsRequest struct {
}
