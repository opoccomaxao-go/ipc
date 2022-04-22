package event

import "sync"

type Pool struct {
	pool sync.Pool
}

const DefaultPayloadSize = 256

func DefaultCommon() interface{} {
	return &Common{
		Payload: make([]byte, DefaultPayloadSize),
	}
}

func NewPool() *Pool {
	return &Pool{
		pool: sync.Pool{
			New: DefaultCommon,
		},
	}
}

func (p *Pool) Acquire() *Common {
	res, ok := p.pool.Get().(*Common)
	if ok {
		return res
	}

	return p.Acquire()
}

func (p *Pool) Release(event *Common) {
	p.pool.Put(event)
}
