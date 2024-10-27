package client

import (
	"context"
	"net/http"
)

func MakeRequest(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	httpsResp, err := GetRequest(ctx, "https://"+url, headers)
	if err != nil {
		httpResp, err := GetRequest(ctx, "http://"+url, headers)
		if err != nil {
			return nil, err
		}
		return httpResp, nil
	}
	return httpsResp, nil
}

func GetRequest(ctx context.Context, url string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header = headers

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
