package integration

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetImages(t *testing.T) {
	t.Parallel()
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
		xCache      string
	}
	tests := []struct {
		name string
		args args
		want results
	}{
		{
			name: "load 200 300 from nginx MISS",
			args: args{
				srv:    "image-previewer:8080",
				width:  200,
				height: 300,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        9407,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "MISS",
			},
		},
		{
			name: "load 200 300 from nginx HIT",
			args: args{
				srv:    "image-previewer:8080",
				width:  200,
				height: 300,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        9407,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "HIT",
			},
		},
		{
			name: "load 300 200 from nginx MISS",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        9690,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "MISS",
			},
		},
		{
			name: "load 1300 200 from nginx MISS",
			args: args{
				srv:    "image-previewer:8080",
				width:  1300,
				height: 200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        22337,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "MISS",
			},
		},
		{
			name: "load 300 1200 from nginx MISS",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        34220,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "MISS",
			},
		},
		{
			name: "load 300 1200 from nginx HIT",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				contentType: "image/jpeg",
				size:        34220,
				format:      "jpeg",
				statusCode:  200,
				xCache:      "HIT",
			},
		},
		{
			name: "wrong server name timeout 408",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "wrong_wrong_server_name/_gopher_original_1024x504.jpg",
			},
			want: results{
				statusCode: 408,
			},
		},
		{
			name: "wrong filename 404",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/wrong_filename",
			},
			want: results{
				statusCode: 404,
			},
		},
		{
			name: "too large 413",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/big.jpg",
			},
			want: results{
				statusCode: 413,
			},
		},
		{
			name: "unsupported media type 415",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: 1200,
				imgURL: "nginx/text.txt",
			},
			want: results{
				statusCode: 415,
			},
		},
		{
			name: "illegal width",
			args: args{
				srv:    "image-previewer:8080",
				width:  -300,
				height: 1200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				statusCode: 400,
			},
		},
		{
			name: "illegal height",
			args: args{
				srv:    "image-previewer:8080",
				width:  300,
				height: -1200,
				imgURL: "nginx/_gopher_original_1024x504.jpg",
			},
			want: results{
				statusCode: 400,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
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
				t.Fatalf("do request: %v", err)
			}
			defer resp.Body.Close()

			log.Printf("response StatusCode: %v", resp.StatusCode)
			require.Equal(t, tt.want.statusCode, resp.StatusCode, "response wrong StatusCode")
			if tt.want.xCache != "" {
				log.Printf("response X-Cache: %s", resp.Header.Get("X-Cache"))
				assert.Equal(t, tt.want.xCache, resp.Header.Get("X-Cache"), "response wrong X-Cache Header")
			}
			if tt.want.contentType != "" {
				log.Printf("response Content-Type: %v", resp.Header.Get("Content-Type"))
				assert.Equal(t, tt.want.contentType, resp.Header.Get("Content-Type"), "response wrong Content-Type")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("read response body: %v", err)
			}
			if tt.want.size != 0 {
				assert.Equal(t, tt.want.size, len(body), "response illegal body size")
			}

			if tt.want.contentType != "image/jpeg" && tt.want.format != "jpeg" && resp.StatusCode == 200 && resp.Header.Get("Content-Type") == "image/jpeg" {
				log.Printf("chack response format: %s", tt.want.format)
				img, imgFormat, err := image.Decode(bytes.NewReader(body))
				if err != nil {
					t.Fatalf("decode image: %v", err)
				}
				assert.Equal(t, tt.args.width, img.Bounds().Dx(), "response illegal image Width")
				assert.Equal(t, tt.args.height, img.Bounds().Dy(), "response illegal image Height")
				assert.Equal(t, tt.want.format, imgFormat, "response illegal image format")
			}
		})
	}
}
