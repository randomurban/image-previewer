package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/randomurban/image-previewer/internal/model"
)

const (
	MaxImageSize = 1000000
)

func MakeRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	resp, err := GetRequest(ctx, "https://"+url, headers)
	if err != nil {
		log.Printf("error on https://%s: %s", url, err)
		resp, err = GetRequest(ctx, "http://"+url, headers)
		if err != nil {
			log.Printf("error on http://%s: %s", url, err)
			return nil, fmt.Errorf("request error: %w", err)
		}
	}
	return resp, nil
}

func GetRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new GET request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	if resp.ContentLength > MaxImageSize {
		return nil, fmt.Errorf("file too big")
	}

	if resp.Header.Get("Content-Type") != "image/jpeg" {
		return nil, fmt.Errorf("file is not image/jpeg")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return &model.ResponseImage{
		Buf: buf,
	}, nil
}
