package client

type Client interface {
	SendToAgent(request interface{}) (interface{}, error)
}
