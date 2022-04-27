package processor

import (
	"github.com/opoccomaxao-go/ipc/channel"
	"github.com/opoccomaxao-go/ipc/event"
)

type Processor struct {
	registry map[uint16]func([]byte)
	_default func(uint16, []byte)
}

// implements interface.
var _ channel.Handler[*event.Common] = (*Processor)(nil)

func (*Processor) discard(uint16, []byte) {}

func (p *Processor) Register(eventType uint16, handler func([]byte)) {
	if handler != nil {
		p.registry[eventType] = handler
	}
}

func (p *Processor) RegisterDefault(handler func(uint16, []byte)) {
	if handler != nil {
		p._default = handler
	}
}

func (p *Processor) Handle(event *event.Common) {
	registered, ok := p.registry[event.Type]
	if ok {
		registered(event.Payload)
	} else {
		p._default(event.Type, event.Payload)
	}
}

func New() *Processor {
	res := &Processor{
		registry: map[uint16]func([]byte){},
	}
	res._default = res.discard

	return res
}
