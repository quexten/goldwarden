package messages

import "encoding/json"

type LockVaultRequest struct {
}

type UnlockVaultRequest struct {
}

type UpdateVaultPINRequest struct {
}

type WipeVaultRequest struct {
}

type GetVaultPINRequest struct {
}

type VaultStatusRequest struct {
}

type VaultStatusResponse struct {
	Locked             bool
	LoggedIn           bool
	PinSet             bool
	NumberOfLogins     int
	NumberOfNotes      int
	LastSynced         int64
	WebsocketConnected bool
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req LockVaultRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, LockVaultRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req UnlockVaultRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, UnlockVaultRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req UpdateVaultPINRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, UpdateVaultPINRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req WipeVaultRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, WipeVaultRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetVaultPINRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetVaultPINRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req VaultStatusRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, VaultStatusRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req VaultStatusResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, VaultStatusResponse{})
}
