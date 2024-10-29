package preview

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/randomurban/image-previewer/internal/http/client"
	"github.com/randomurban/image-previewer/internal/service"
	"github.com/randomurban/image-previewer/internal/storage"
)

const (
	HTTPClientTimeout = 5 * time.Second
)

var _ service.Previewer = (*Preview)(nil)

type Preview struct {
	store storage.Cacher
}

func NewPreviewService(store storage.Cacher) service.Previewer {
	return &Preview{
		store: store,
	}
}

func (s *Preview) PreviewImage(width int, height int, url string, header http.Header) ([]byte, error) {
	clientCtx, clientCancel := context.WithTimeout(context.Background(), HTTPClientTimeout)
	defer clientCancel()

	name := sha256.Sum256([]byte(fmt.Sprintf("%v_%v_%v", width, height, url)))
	key := base32.StdEncoding.EncodeToString(name[:])
	fromCache, err := s.store.Download(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get image from cache: %w", err)
	}
	if fromCache != nil {
		log.Printf("get image from cache")
		return fromCache, nil
	}

	resp, err := client.MakeRequest(clientCtx, url, header)
	if err != nil {
		return nil, err
	}

	originalImg, err := jpeg.Decode(bytes.NewReader(resp.Buf))
	if err != nil {
		return nil, err
	}
	// Lanczos Linear NearestNeighbor
	resizedImg := imaging.Fill(originalImg, width, height, imaging.Center, imaging.Linear)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, resizedImg, nil)
	if err != nil {
		return nil, err
	}

	err = s.store.Upload(key, buf.Bytes())

	return buf.Bytes(), nil
}
