package channel

import (
	"net"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/transport"
	"github.com/pkg/errors"
)

type Client struct {
	*Channel
	config ClientConfig
}

type ClientConfig struct {
	Handler Handler[*event.Common] // Handler handles incoming events
	Address string                 // PublicAddr is used for listen (Server) / connect (Client)
}

func Connect(config ClientConfig) (*Client, error) {
	if config.Address == "" {
		return nil, errors.WithMessage(ErrNoParam, "Address")
	}

	if config.Handler == nil {
		return nil, errors.WithMessage(ErrNoParam, "Handler")
	}

	res := &Client{
		config: config,
	}

	transport, err := res.connect()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res.Channel = &Channel{
		Connect:   res.connect,
		Transport: transport,
		TryCount:  100,
	}

	return res, nil
}

func (c *Client) connect() (transport.Transport, error) {
	conn, err := net.Dial("tcp", c.config.Address)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return transport.NewSocket(conn), nil
}
