package transport

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestStream(t *testing.T) {
	t.Parallel()

	buffer := bytes.Buffer{}
	transport := NewStream(&buffer)

	suite.Run(t, NewTestSuite(transport, transport))
}
