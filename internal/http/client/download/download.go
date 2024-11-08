package download

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/randomurban/image-previewer/internal/model"
)

type Client struct {
	MaxImageSize int64
}

func NewClient(maxImageSize int64) *Client {
	return &Client{
		MaxImageSize: maxImageSize,
	}
}

func (c *Client) MakeRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	resp, err := c.GetRequest(ctx, "https://"+url, headers)
	if err != nil {
		log.Printf("error on https://%s: %s", url, err)
		resp, err = c.GetRequest(ctx, "http://"+url, headers)
		if err != nil {
			log.Printf("error on http://%s: %s", url, err)
			return nil, fmt.Errorf("request error: %w", err)
		}
	}
	return resp, nil
}

func (c *Client) GetRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new GET request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("error closing body: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	if resp.ContentLength > c.MaxImageSize {
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
