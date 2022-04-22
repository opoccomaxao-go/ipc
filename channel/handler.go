package channel

type Handler[T any] interface {
	Handle(T)
}

type ChannelHandler interface {
	Handle()
}
