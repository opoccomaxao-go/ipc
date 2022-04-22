package transport

import (
	"github.com/opoccomaxao-go/ipc/event"
)

type Transport interface {
	Write(event event.Event) error
	Read(eventRef event.Event) error
	// Close should close transport if it has such option.
	Close() error
}
