package service

import "net/http"

type Previewer interface {
	PreviewImage(width int, height int, url string, header http.Header) ([]byte, error)
}
