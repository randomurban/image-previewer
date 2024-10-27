package model

import (
	"image"
	"io"
	"net/http"
)

type Preview struct {
	OriginalUrl string
	Buffer      io.Reader
	Header      http.Header
	Image       image.Image
	Width       int
	Height      int
	Filename    string
}

type ResponseImage struct {
	Buf    []byte
	Header *http.Header
}
