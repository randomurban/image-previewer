package download

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/randomurban/image-previewer/internal/model"
)

type Client struct {
	MaxImageSize int64
	ClientTimout time.Duration
}

func NewClient(maxImageSize int64, timeout time.Duration) *Client {
	return &Client{
		MaxImageSize: maxImageSize,
		ClientTimout: timeout,
	}
}

func (c *Client) MakeRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	log.Printf("trying https://%s", url)
	resp, err := c.GetRequest(ctx, "https://"+url, headers)
	if err != nil {
		log.Printf("response from https: %s", err)
		log.Printf("trying http://%s", url)
		resp, err = c.GetRequest(ctx, "http://"+url, headers)
		if err != nil {
			log.Printf("response from http: %s", err)
			return nil, fmt.Errorf("request http: %w", err)
		}
	}
	return resp, nil
}

func (c *Client) GetRequest(ctx context.Context, url string, headers http.Header) (*model.ResponseImage, error) {
	reqCtx, reqCancel := context.WithTimeout(ctx, c.ClientTimout)
	defer reqCancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("new GET request: %w: %w", model.ErrBadGateway, err)
	}
	req = req.WithContext(reqCtx)
	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return nil, fmt.Errorf("request timeout: %w", model.ErrTimeout)
		}
		return nil, fmt.Errorf("do request: %w: %w", model.ErrRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, model.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d", resp.StatusCode)
	}

	if resp.ContentLength > c.MaxImageSize {
		return nil, fmt.Errorf("file too big: %w", model.ErrTooLarge)
	}

	if resp.Header.Get("Content-Type") != "image/jpeg" {
		return nil, fmt.Errorf("file is not image/jpeg")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w: %w", model.ErrBadGateway, err)
	}
	return &model.ResponseImage{
		Buf: buf,
	}, nil
}
