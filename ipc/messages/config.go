package messages

import "encoding/json"

type SetApiURLRequest struct {
	Value string
}

type SetIdentityURLRequest struct {
	Value string
}

type SetNotificationsURLRequest struct {
	Value string
}

type SetClientIDRequest struct {
	Value string
}

type SetClientSecretRequest struct {
	Value string
}

type GetRuntimeConfigRequest struct{}

type GetRuntimeConfigResponse struct {
	UseMemguard          bool
	SSHAgentSocketPath   string
	GoldwardenSocketPath string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SetApiURLRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SetApiURLRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SetIdentityURLRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SetIdentityURLRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SetNotificationsURLRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SetNotificationsURLRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetRuntimeConfigRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetRuntimeConfigRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetRuntimeConfigResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetRuntimeConfigResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SetClientIDRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SetClientIDRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req SetClientSecretRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, SetClientSecretRequest{})
}
