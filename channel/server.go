package channel

import (
	"net"

	"github.com/opoccomaxao-go/ipc/transport"
	"github.com/pkg/errors"
)

type Server struct {
	listener net.Listener
	config   ServerConfig
}

type ServerConfig struct {
	Address string            // Address with port to listen.
	Handler Handler[*Channel] // Handler handles every connected client.
}

func NewServer(config ServerConfig) (*Server, error) {
	if config.Handler == nil {
		return nil, errors.WithMessage(ErrNoParam, "Handler")
	}

	return &Server{
		config: config,
	}, nil
}

func (s *Server) Listen() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return errors.WithStack(err)
	}

	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.WithStack(err)
		}

		transport := transport.NewSocket(conn)

		go s.config.Handler.Handle(&Channel{
			Transport: transport,
			Connect:   nil,
			TryCount:  1,
		})
	}
}

func (s *Server) Close() error {
	return errors.WithStack(s.listener.Close())
}
