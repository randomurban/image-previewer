package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

const (
	EnvPrefix = "PREVIEWER_"

	HttpAddr        = "HTTP_ADDR"
	HttpAddrDefault = "localhost:8080"

	CacheCap        = "CACHE_CAP"
	CacheCapDefault = 10
)

type Config struct {
	HttpAddr string
	CacheCap int
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
		HttpAddr: env(EnvPrefix+HttpAddr, HttpAddrDefault),
		CacheCap: envInt(EnvPrefix+CacheCap, CacheCapDefault),
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
