package channel

import (
	"sync"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/transport"
	"github.com/pkg/errors"
)

type ConnectFn func() (transport.Transport, error)

type Channel struct {
	Transport transport.Transport
	Connect   ConnectFn
	TryCount  int
	closed    bool
	stateMu   sync.Mutex
	sendMu    sync.Mutex
}

func (c *Channel) Send(event *event.Common) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	for i := c.TryCount; i > 0; i-- {
		err := errors.WithStack(c.Transport.Write(event))
		if err == nil {
			return nil
		}

		if errors.Is(err, transport.ErrClosed) {
			err := c.reconnect()
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return errors.WithStack(ErrConnectionFailed)
}

func (c *Channel) Serve(handler Handler[*event.Common]) error {
	var buffer event.Common

	for {
		err := errors.WithStack(c.Transport.Read(&buffer))
		if err == nil {
			handler.Handle(&buffer)
		} else if errors.Is(err, transport.ErrClosed) {
			err := c.reconnect()
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
}

func (c *Channel) reconnect() error {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()

	if c.closed || c.Connect == nil {
		return errors.WithStack(transport.ErrClosed)
	}

	transport, err := c.Connect()
	if err != nil {
		return errors.WithStack(err)
	}

	c.Transport = transport

	return nil
}

func (c *Channel) Close() error {
	c.stateMu.Lock()
	c.closed = true
	c.stateMu.Unlock()

	c.sendMu.Lock()
	defer c.sendMu.Unlock()

	return errors.WithStack(c.Transport.Close())
}
