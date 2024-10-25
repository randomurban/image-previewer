package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/randomurban/image-previewer/internal/config"
	"github.com/randomurban/image-previewer/internal/http/server"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", ".env", "path to config file")
	flag.Parse()
}

func main() {
	log.Printf("load config from: %v", configPath)
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer cancel()

	httpServer := server.NewHTTPServer()
	go func() {
		log.Printf("server started at: http://%s", cfg.HttpAddr)
		err := httpServer.Start(cfg.HttpAddr)
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %s", err)
		}
		cancel()
	}()

	<-ctx.Done()
	log.Printf("server is stopping...")

	timeOut, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer timeOutCancel()

	if err := httpServer.Stop(timeOut); err != nil {
		log.Printf("server stop: %s", err)
	}
	log.Printf("server is stopped")
}
