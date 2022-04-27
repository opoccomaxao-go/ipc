package processor

import (
	"math"
	"testing"

	"github.com/opoccomaxao-go/ipc/event"
	"github.com/stretchr/testify/require"
)

//nolint:tparallel
func TestProcessor(t *testing.T) {
	t.Parallel()

	var res [math.MaxUint16 + 1]bool

	proc := New()

	testCases := []struct {
		desc    string
		prepare func()
	}{
		{
			desc: "Default",
			prepare: func() {
				proc.RegisterDefault(func(u uint16, b []byte) { res[u] = len(b) > 0 })
			},
		},
		{
			desc: "Registered",
			prepare: func() {
				for i := 0; i <= math.MaxUint16; i++ {
					i := i
					proc.Register(uint16(i), func(b []byte) { res[i] = len(b) > 0 })
				}
			},
		},
	}

	//nolint:paralleltest
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			tC.prepare()

			for i := range res {
				res[i] = false
			}

			event := event.Common{
				Type:    0,
				Payload: []byte{0},
			}

			for i := 0; i <= math.MaxUint16; i++ {
				event.Type = uint16(i)
				proc.Handle(&event)
			}

			for i, v := range res {
				require.True(t, v, i)
			}
		})
	}
}
