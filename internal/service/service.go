package service

import (
	"net/http"

	"github.com/randomurban/image-previewer/internal/model"
)

type Previewer interface {
	PreviewImage(width int, height int, url string, header http.Header) (*model.ResponseImage, error)
}
