package channel

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/opoccomaxao-go/ipc/mock"
	"github.com/opoccomaxao-go/ipc/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientServer(t *testing.T) {
	t.Parallel()

	serverEventsToSend := []*event.Common{
		{Type: 1, Payload: mock.Bytes(10000)},
		{Type: 2, Payload: mock.Bytes(20000)},
		{Type: 3, Payload: mock.Bytes(30000)},
		{Type: 4, Payload: mock.Bytes(40000)},
	}

	clientEventsToSend := []*event.Common{
		{Type: 5000, Payload: mock.Bytes(1000)},
		{Type: 6000, Payload: mock.Bytes(2000)},
		{Type: 7000, Payload: mock.Bytes(3000)},
		{Type: 8000, Payload: mock.Bytes(4000)},
	}

	var serverRecvWG sync.WaitGroup
	serverRecvWG.Add(len(clientEventsToSend))

	var clientRecvWG sync.WaitGroup
	clientRecvWG.Add(len(serverEventsToSend))

	serverReceiver := CollectorHandler[*event.Common]{
		OnHandle: serverRecvWG.Done,
	}
	clientReceiver := CollectorHandler[*event.Common]{
		OnHandle: clientRecvWG.Done,
	}

	var sendData = func(data []*event.Common, channel *Channel) {
		for _, e := range data {
			err := channel.Send(e)

			assert.NoError(t, err)
		}
	}

	var server *Server
	server, err := NewServer(ServerConfig{
		Address: ":10001",
		Handler: &CommonHandler[*Channel]{
			Func: func(connection *Channel) {
				go sendData(serverEventsToSend, connection)
				go func() {
					err := connection.Serve(&serverReceiver)
					if err != nil {
						require.ErrorIs(t, err, transport.ErrClosed)
					}
				}()
			},
		},
	})

	require.NoError(t, err)

	go func() {
		err := server.Listen()

		require.ErrorIs(t, err, net.ErrClosed)
	}()

	// wait for server start. doesn't work without this
	time.Sleep(time.Second)

	client, err := Dial(ClientConfig{
		Address: "127.0.0.1:10001",
		Handler: &clientReceiver,
	})
	require.NoError(t, err)

	go sendData(clientEventsToSend, client.Channel)
	go func() {
		err = client.Serve(&clientReceiver)
		if err != nil {
			require.ErrorIs(t, err, transport.ErrClosed)
		}
	}()

	serverRecvWG.Wait()
	clientRecvWG.Wait()
	require.NoError(t, server.Close())
	require.NoError(t, client.Close())

	assert.Equal(t, serverEventsToSend, clientReceiver.Collect())
	assert.Equal(t, clientEventsToSend, serverReceiver.Collect())
}
