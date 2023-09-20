package messages

import "encoding/json"

type ActionResponse struct {
	Success bool
	Message string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req ActionResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, ActionResponse{})
}
