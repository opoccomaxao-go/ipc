package channel

import (
	"sync"
)

type Handler[T any] interface {
	Handle(T)
}

type HandlerFunc[T any] func(T)

type CommonHandler[T any] struct {
	Func HandlerFunc[T]
}

func (h *CommonHandler[T]) Handle(value T) {
	h.Func(value)
}

// CollectorHandler for test purposes.
type CollectorHandler[T interface {
	Copy() T
}] struct {
	OnHandle func()

	data []T
	mu   sync.Mutex
}

func (h *CollectorHandler[T]) Handle(value T) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.data = append(h.data, value.Copy())

	if h.OnHandle != nil {
		h.OnHandle()
	}
}

func (h *CollectorHandler[T]) Collect() []T {
	h.mu.Lock()
	defer h.mu.Unlock()

	res := h.data
	h.data = nil

	return res
}
