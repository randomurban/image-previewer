package server

import (
	"bytes"
	"context"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/randomurban/image-previewer/internal/service"
)

type Server struct {
	srv *http.Server
	mux *http.ServeMux
}

func NewHTTPServer() *Server {
	return &Server{srv: &http.Server{}, mux: http.NewServeMux()}
}

func (s *Server) Start(addr string) error {
	s.srv = &http.Server{Addr: addr, Handler: s.mux}

	s.mux.HandleFunc("GET /fill/{width}/{height}/{img...}", s.FillHandle)

	err := s.srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) FillHandle(w http.ResponseWriter, r *http.Request) {
	width, err := strconv.Atoi(r.PathValue("width"))
	if err != nil {
		log.Printf("width: %v", err)
		http.Error(w, "Bad request: wrong width in url", http.StatusBadRequest)
		return
	}
	log.Printf("width: %v", width)

	height, err := strconv.Atoi(r.PathValue("height"))
	if err != nil {
		log.Printf("height: %v", err)
		http.Error(w, "Bad request: wrong height in url", http.StatusBadRequest)
		return
	}
	log.Printf("height: %v", height)

	url := r.PathValue("img")
	log.Printf("image url: %v", url)

	img, respHeader, err := service.PreviewImage(width, height, url, r.Header)
	if err != nil {
		log.Printf("preview image: %v", err)
		http.Error(w, "preview: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for key, values := range respHeader {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		log.Printf("encode: %v", err)
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Printf("failed write: %s", err)
		return
	}
}
