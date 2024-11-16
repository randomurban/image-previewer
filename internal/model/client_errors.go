package model

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrRequest    = errors.New("url error")
	ErrTooLarge   = errors.New("max image size exceeds limit")
	ErrBadGateway = errors.New("bad gateway")
	ErrTimeout    = errors.New("timeout")
)
