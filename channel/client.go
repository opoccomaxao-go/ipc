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
	Handler   Handler[*event.Common] // Handler handles incoming events.
	Address   string                 // Address is used to connect.
	Reconnect bool                   // Reconnect if reconnection needed.
}

func Dial(config ClientConfig) (*Client, error) {
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
		Transport: transport,
		TryCount:  100,
	}

	if config.Reconnect {
		res.Channel.Connect = res.connect
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
