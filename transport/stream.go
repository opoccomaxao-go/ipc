package transport

import (
	"io"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/pkg/errors"
)

type Stream struct {
	source io.ReadWriter
}

// implements interface.
var _ Transport = (*Stream)(nil)

func (s *Stream) Write(event event.Event) error {
	return errors.WithStack(event.WriteBinary(s.source))
}

func (s *Stream) Read(event event.Event) error {
	return errors.WithStack(event.ReadBinary(s.source))
}

func (s *Stream) Close() error {
	if c, ok := s.source.(io.Closer); ok {
		return errors.WithStack(c.Close())
	}

	return nil
}

func NewStream(stream io.ReadWriter) Transport {
	return &Stream{
		source: stream,
	}
}
