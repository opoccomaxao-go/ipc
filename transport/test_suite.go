package transport

import (
	"github.com/opoccomaxao-go/ipc/event"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	input  Transport
	output Transport
	suite.Suite
}

// NewTestSuite creates test suite for Transport.
// Before tests create and connect output to input.
// See example:
//  transport/stream_test.go
//  transport/tcp_test.go.
func NewTestSuite(input Transport, output Transport) *TestSuite {
	return &TestSuite{
		input:  input,
		output: output,
	}
}

func (suite *TestSuite) TestAMultipleEvents() {
	//nolint:gomnd // test data
	batches := [][]event.Common{
		{
			{Type: 1, Payload: []byte("test")},
		},
		{
			{Type: 2, Payload: []byte("test 2")},
			{Type: 3, Payload: []byte("test 3")},
		},
		{
			{Type: 4, Payload: []byte("test 4")},
			{Type: 5, Payload: []byte("test 5")},
			{Type: 6, Payload: []byte("test 6")},
		},
		{
			{Type: 7, Payload: make([]byte, 10000)},
			{Type: 8, Payload: []byte("test 8")},
			{Type: 9, Payload: make([]byte, 10000)},
			{Type: 10, Payload: []byte("test 8")},
		},
	}

	checkList := [][2]Transport{
		{suite.input, suite.output},
		{suite.output, suite.input},
	}

	for _, list := range checkList {
		input := list[0]
		output := list[1]

		for batchID, batch := range batches {
			for i := range batch {
				err := output.Write(&batch[i])
				suite.Require().NoError(err, batchID)
			}

			res := make([]event.Common, len(batch))
			for i := range res {
				suite.Require().NoError(input.Read(&res[i]), i, batchID)
			}

			suite.Require().Equal(batch, res, batchID)
		}
	}
}

// TestZClosed should be final test.
func (suite *TestSuite) TestZClosed() {
	var event event.Common

	list := []Transport{suite.input, suite.output}

	for _, transport := range list {
		_ = transport.Close()
	}

	for _, transport := range list {
		err := transport.Write(&event)
		if err != nil {
			suite.Require().ErrorIs(err, ErrClosed)
		}
	}

	for _, transport := range list {
		err := transport.Read(&event)
		if err != nil {
			suite.Require().ErrorIs(err, ErrClosed)
		}
	}
}
