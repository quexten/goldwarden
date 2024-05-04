package messages

import "encoding/json"

type GetBiometricsKeyRequest struct {
}

type GetBiometricsKeyResponse struct {
	Key string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetBiometricsKeyRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetBiometricsKeyRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetBiometricsKeyResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetBiometricsKeyResponse{})
}
