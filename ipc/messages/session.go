package messages

import "encoding/json"

type SessionAuthRequest struct {
	Token string
}

type SessionAuthResponse struct {
	Verified bool
}

type PinentryRegistrationRequest struct {
}

type PinentryRegistrationResponse struct {
	Success bool
}

type PinentryPinRequest struct {
	Message string
}

type PinentryPinResponse struct {
	Pin string
}

type PinentryApprovalRequest struct {
	Message string
}

type PinentryApprovalResponse struct {
	Approved bool
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

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryRegistrationRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryRegistrationRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryRegistrationResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryRegistrationResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryPinRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryPinRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryPinResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryPinResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryApprovalRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryApprovalRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req PinentryApprovalResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, PinentryApprovalResponse{})
}
