package messages

import "encoding/json"

type SessionAuthRequest struct {
	Token string
}

type SessionAuthResponse struct {
	Verified bool
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SessionAuthRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SessionAuthRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SessionAuthResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SessionAuthResponse{})
}
