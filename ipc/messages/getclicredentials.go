package messages

import "encoding/json"

type GetCLICredentialsRequest struct {
	ApplicationName string
}

type GetCLICredentialsResponse struct {
	Env map[string]string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetCLICredentialsRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetCLICredentialsRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetCLICredentialsResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetCLICredentialsResponse{})
}
