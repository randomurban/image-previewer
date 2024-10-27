package service

import (
	"bytes"
	"context"
	"image/jpeg"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/randomurban/image-previewer/internal/http/client"
)

const (
	HTTPClientTimeout = 5 * time.Second
)

func PreviewImage(width int, height int, url string, header http.Header) ([]byte, error) {
	clientCtx, clientCancel := context.WithTimeout(context.Background(), HTTPClientTimeout)
	defer clientCancel()

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
	return buf.Bytes(), nil
}
