package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/randomurban/image-previewer/internal/http/server"
)

func main() {
	addr := "localhost:8080"

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	httpServer := server.NewHTTPServer()
	go func() {
		log.Printf("server started at: http://%s", addr)
		err := httpServer.Start(addr)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %s", err)
		}
		cancel()
	}()

	<-ctx.Done()
	log.Printf("server is stoping...")

	timeOut, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer timeOutCancel()

	if err := httpServer.Stop(timeOut); err != nil {
		log.Printf("server stop: %s", err)
	}
	log.Printf("server stoped")
}
