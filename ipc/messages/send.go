package messages

import "encoding/json"

type GetSendRequest struct {
	Name string
	Text string
}

type GetSendResponse struct {
	Found bool
	Text  string
}

type CreateSendRequest struct {
	Name string
	Text string
}

type CreateSendResponse struct {
	URL string
}

type ListSendsRequest struct {
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetSendRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetSendRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetSendResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetSendResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req CreateSendRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, CreateSendRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req CreateSendResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, CreateSendResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req ListSendsRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, ListSendsRequest{})
}
