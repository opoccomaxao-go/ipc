package channel

import "errors"

var (
	ErrConnectionFailed = errors.New("connection failed")
	ErrNoParam          = errors.New("no param")
)
