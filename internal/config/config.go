package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	EnvPrefix = "PREVIEWER_"

	HTTPAddr        = "HTTP_ADDR"
	HTTPAddrDefault = "localhost:8080"

	CacheDir        = "CACHE_DIR"
	CacheDirDefault = "./cache_store"

	CacheCap        = "CACHE_CAP"
	CacheCapDefault = 10

	HTTPClientTimeout        = "HTTP_CLIENT_TIMEOUT"
	HTTPClientTimeoutDefault = 5

	HTTPServerReadTimeout        = "HTTP_SERVER_READ_TIMEOUT"
	HTTPServerReadTimeoutDefault = 10

	HTTPServerWriteTimeout        = "HTTP_SERVER_WRITE_TIMEOUT"
	HTTPServerWriteTimeoutDefault = 10

	MaxImageSize        = "MAX_IMAGE_SIZE"
	MaxImageSizeDefault = 1000000
)

type Config struct {
	HTTPAddr               string
	CacheDir               string
	CacheCap               int
	HTTPClientTimeout      time.Duration
	HTTPServerReadTimeout  time.Duration
	HTTPServerWriteTimeout time.Duration
	MaxImageSize           int64
}

func New(configPath string) (*Config, error) {
	envPath := configPath
	if configPath == "" {
		envPath = ".env"
	}
	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("load %s: %w", envPath, err)
	}
	res := &Config{
		HTTPAddr:               env(EnvPrefix+HTTPAddr, HTTPAddrDefault),
		CacheDir:               env(EnvPrefix+CacheDir, CacheDirDefault),
		CacheCap:               envInt(EnvPrefix+CacheCap, CacheCapDefault),
		HTTPClientTimeout:      envSecond(EnvPrefix+HTTPClientTimeout, HTTPClientTimeoutDefault),
		HTTPServerReadTimeout:  envSecond(EnvPrefix+HTTPServerReadTimeout, HTTPServerReadTimeoutDefault),
		HTTPServerWriteTimeout: envSecond(EnvPrefix+HTTPServerWriteTimeout, HTTPServerWriteTimeoutDefault),
		MaxImageSize:           envInt64(EnvPrefix+MaxImageSize, MaxImageSizeDefault),
	}
	return res, nil
}

func env(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if ok {
		return val
	}
	return defaultVal
}

func envInt(key string, defaultVal int) int {
	val := env(key, "")
	res, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return res
}

func envSecond(key string, defaultVal time.Duration) time.Duration {
	val := env(key, "")
	resint, err := strconv.Atoi(val)
	res := time.Duration(resint) * time.Second
	if err != nil {
		return defaultVal
	}
	return res
}

func envInt64(key string, defaultVal int64) int64 {
	val := env(key, "")
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultVal
	}
	return int64(res)
}
