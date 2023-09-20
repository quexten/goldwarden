package messages

import "encoding/json"

type GetNotesRequest struct {
	Name string
}

type GetNoteResponse struct {
	Found  bool
	Result DecryptedNoteCipher
}

type GetNotesResponse struct {
	Found  bool
	Result []DecryptedNoteCipher
}

type DecryptedNoteCipher struct {
	Name     string
	Contents string
}

func init() {
	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetNotesRequest
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetNotesRequest{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetNoteResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetNoteResponse{})

	registerPayloadParser(func(payload []byte) (interface{}, error) {
		var req GetNotesResponse
		err := json.Unmarshal(payload, &req)
		if err != nil {
			panic("Unmarshal: " + err.Error())
		}
		return req, nil
	}, GetNotesResponse{})
}
