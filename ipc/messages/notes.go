package messages

import "encoding/json"

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
	TwoFactorCode string
}

type AddLoginRequest struct {
	Name string
	UUID string
}

type AddLoginResponse struct {
	Name string
	UUID string
}

type ListLoginsRequest struct {
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetLoginRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetLoginRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetLoginResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetLoginResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetLoginsResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetLoginsResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req AddLoginRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, AddLoginRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req AddLoginResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, AddLoginResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req ListLoginsRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, ListLoginsRequest{})
}
