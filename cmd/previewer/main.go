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
	"github.com/randomurban/image-previewer/internal/http/client/download"
	"github.com/randomurban/image-previewer/internal/http/server/handle/fill"
	"github.com/randomurban/image-previewer/internal/service/preview"
	"github.com/randomurban/image-previewer/internal/storage/filestorage"
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

	cacheStore := filestorage.NewStorage(cfg.CacheDir, cfg.CacheCap)
	err = cacheStore.Init()
	if err != nil {
		log.Fatal("cache init: ", err)
	}
	defer cancel()

	downloader := download.NewClient(cfg.MaxImageSize, cfg.HTTPClientTimeout)
	previewer := preview.NewPreviewService(cacheStore, downloader, cfg.HTTPClientTimeout)

	fillHandler := fill.NewHandle(previewer)

	router := http.NewServeMux()
	router.HandleFunc("GET /fill/{width}/{height}/{img...}", fillHandler.FillHandle)

	httpServer := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServerReadTimeout,
		WriteTimeout: cfg.HTTPServerWriteTimeout,
	}

	go func() {
		log.Printf("server started at: http://%s", cfg.HTTPAddr)
		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server: %s", err)
		}
		cancel()
	}()

	<-ctx.Done()
	log.Printf("server is stopping...")

	timeOut, timeOutCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer timeOutCancel()

	if err := httpServer.Shutdown(timeOut); err != nil {
		log.Printf("server stop: %s", err)
	}
	log.Printf("server is stopped")
}
