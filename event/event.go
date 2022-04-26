package event

import (
	"io"

	"github.com/pkg/errors"
)

type Event interface {
	WriteBinary(writer io.Writer) error
	ReadBinary(reader io.Reader) error
}

type Common struct {
	Type    uint16
	Payload []byte
}

type eventHeader [4]byte

//nolint:gomnd
func (e *Common) WriteBinary(writer io.Writer) error {
	size := len(e.Payload)

	_, err := writer.Write([]byte{
		byte(e.Type & 0xff),
		byte(e.Type >> 8 & 0xff),
		byte(size & 0xff),
		byte(size >> 8 & 0xff),
	})
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = writer.Write(e.Payload)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//nolint:gomnd
func (e *Common) ReadBinary(reader io.Reader) error {
	var header eventHeader

	_, err := reader.Read(header[:])
	if err != nil {
		return errors.WithStack(err)
	}

	e.Type = uint16(header[0]) + uint16(header[1])<<8
	size := int(header[2]) + int(header[3])<<8

	if cap(e.Payload) >= size {
		e.Payload = e.Payload[0:size]
	} else {
		e.Payload = make([]byte, size)
	}

	_, err = reader.Read(e.Payload)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (e *Common) Copy() *Common {
	res := &Common{
		Type:    e.Type,
		Payload: make([]byte, len(e.Payload)),
	}

	copy(res.Payload, e.Payload)

	return res
}
