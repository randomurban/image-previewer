package integration

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetImages(t *testing.T) {
	type args struct {
		srv           string
		width, height int
		imgURL        string
	}
	type results struct {
		contentType string
		size        int
		format      string
		statusCode  int
	}
	tests := []struct {
		name string
		args args
		want results
	}{
		{
			name: "load 200 300 from nginx",
			args: args{
				srv:    "image-previewer:8080",
				width:  200,
				height: 300,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        18814,
				format:      "jpeg",
				statusCode:  200,
			},
		},
		{
			name: "load 300 200 from nginx",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        19380,
				format:      "jpeg",
				statusCode:  200,
			},
		},
		{
			name: "load 1300 200 from nginx",
			args: args{
				srv:    "image-previewer:8080",
				width:  1300,
				height: 200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        44674,
				format:      "jpeg",
				statusCode:  200,
			},
		},
		{
			name: "load 300 1200 from nginx",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        68440,
				format:      "jpeg",
				statusCode:  200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			url := fmt.Sprintf("http://%s/fill/%v/%v/%s",
				tt.args.srv, tt.args.width, tt.args.height, tt.args.imgURL)

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				t.Fatalf("new request: %v", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("do request: %v", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, tt.want.statusCode, resp.StatusCode)
			assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"))

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("read response body: %v", err)
			}
			assert.Equal(t, tt.want.size, len(body), "illegal body size")

			img, imgFormat, err := image.Decode(bytes.NewReader(body))
			if err != nil {
				t.Fatalf("decode image: %v", err)
			}
			assert.Equal(t, tt.args.width, img.Bounds().Dx(), "illegal image Dx")
			assert.Equal(t, tt.args.height, img.Bounds().Dy(), "illegal image Dy")
			assert.Equal(t, tt.want.format, imgFormat, "illegal image format")
		})
	}
}
