package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

	s.mux.HandleFunc("GET /fill/{img...}", s.FillHandle)

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
	img := r.PathValue("img")
	fmt.Printf("got %s\n", img)
	_, err := fmt.Fprintf(w, "img: %s\n", img)
	if err != nil {
		log.Printf("failed write: %s", err)
	}
}
