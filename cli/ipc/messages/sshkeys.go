package messages

import "encoding/json"

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

type ImportSSHKeyRequest struct {
	Key  string
	Name string
}

type ImportSSHKeyResponse struct {
	Success  bool
	ErrorMsg string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req CreateSSHKeyRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, CreateSSHKeyRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req CreateSSHKeyResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, CreateSSHKeyResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetSSHKeysRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetSSHKeysRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetSSHKeysResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetSSHKeysResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req ImportSSHKeyRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, ImportSSHKeyRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req ImportSSHKeyResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, ImportSSHKeyResponse{})
}
