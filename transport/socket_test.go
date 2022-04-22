package transport

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestSocket(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	var servConn, clientConn net.Conn

	go func() {
		var err error
		clientConn, err = net.Dial("tcp", listener.Addr().String())
		require.NoError(t, err)
	}()

	servConn, err = listener.Accept()
	require.NoError(t, err)

	servTransport := NewSocket(servConn)
	clientTransport := NewSocket(clientConn)

	suite.Run(t, NewTestSuite(servTransport, clientTransport))
}
