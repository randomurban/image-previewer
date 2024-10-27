package service

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"time"

	"github.com/disintegration/imaging"
	"github.com/randomurban/image-previewer/internal/http/client"
)

const (
	HttpClientTimeout = 5 * time.Second
	MaxImageSize      = 1000000
)

func PreviewImage(width int, height int, url string, header http.Header) (image.Image, http.Header, error) {
	originalImg, respHeader, err := downloadImage(url, header)
	if err != nil {
		return nil, nil, err
	}

	// Lanczos Linear NearestNeighbor
	resizedImg := imaging.Fill(originalImg, width, height, imaging.Center, imaging.Linear)

	return resizedImg, respHeader, nil
}

func downloadImage(url string, header http.Header) (image.Image, http.Header, error) {
	clientCtx, clientCancel := context.WithTimeout(context.Background(), HttpClientTimeout)
	defer clientCancel()

	resp, err := client.MakeRequest(clientCtx, url, header)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf(resp.Status)
	}

	if resp.ContentLength > MaxImageSize {
		return nil, nil, fmt.Errorf("file too big")
	}

	if resp.Header.Get("Content-Type") != "image/jpeg" {
		return nil, nil, fmt.Errorf("file is not image/jpeg")
	}

	originalImg, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return originalImg, resp.Header, nil
}
