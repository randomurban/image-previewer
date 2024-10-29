package client

import (
	"context"
	"net/http"

	"github.com/randomurban/image-previewer/internal/model"
)

type Downloader interface {
	MakeRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error)
	GetRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error)
}
