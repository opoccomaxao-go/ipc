package transport

import (
	"net"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/pkg/errors"
)

type Socket struct {
	socket net.Conn
}

// implements interface.
var _ Transport = &Socket{}

func (*Socket) checkError(err error) error {
	if errors.Is(err, net.ErrClosed) {
		return ErrClosed
	}

	return err
}

func (s *Socket) Write(event event.Event) error {
	return errors.WithStack(
		s.checkError(
			event.WriteBinary(s.socket),
		),
	)
}

func (s *Socket) Read(event event.Event) error {
	return errors.WithStack(
		s.checkError(
			event.ReadBinary(s.socket),
		),
	)
}

func (s *Socket) Close() error {
	return errors.WithStack(s.socket.Close())
}

func NewSocket(socket net.Conn) Transport {
	return &Socket{
		socket: socket,
	}
}
