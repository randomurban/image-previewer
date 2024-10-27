package server

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/randomurban/image-previewer/internal/service"
)

type Server struct {
	srv *http.Server
}

func NewHTTPServer(addr string) *Server {
	router := http.NewServeMux()
	router.HandleFunc("GET /fill/{width}/{height}/{img...}", FillHandle)

	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	err := s.srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func FillHandle(w http.ResponseWriter, r *http.Request) {
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

	imgBuf, err := service.PreviewImage(width, height, url, r.Header)
	if err != nil {
		log.Printf("preview image: %v", err)
		http.Error(w, "preview: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(imgBuf)
	if err != nil {
		log.Printf("encode: %v", err)
		http.Error(w, "encode error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(imgBuf)))
	_, err = w.Write(imgBuf)
	if err != nil {
		log.Printf("failed write: %s", err)
		return
	}
}
