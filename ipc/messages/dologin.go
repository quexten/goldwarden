package messages

import "encoding/json"

type DoLoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	Passwordless bool   `json:"passwordless"`
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req DoLoginRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, DoLoginRequest{})
}
