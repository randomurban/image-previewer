package model

import "errors"

var (
	ErrNotFound             = errors.New("not found")
	ErrTooLarge             = errors.New("max image size exceeds limit")
	ErrBadGateway           = errors.New("bad gateway")
	ErrTimeout              = errors.New("timeout")
	ErrUnsupportedMediaType = errors.New("unsupported media type")
)
