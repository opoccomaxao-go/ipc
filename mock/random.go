package mock

import (
	"math/rand"
	"time"
)

func Bytes(length int) []byte {
	res := make([]byte, length)

	rand.Seed(time.Now().Unix())

	//nolint:gosec // for test purposes
	_, _ = rand.Read(res)

	return res
}
