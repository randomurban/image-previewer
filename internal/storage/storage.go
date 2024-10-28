package storage

import "github.com/randomurban/image-previewer/internal/cache"

type Storage struct {
	cache    cache.Cache
	capacity int
	path     string
}

func NewStorage(path string, capacity int) *Storage {
	return &Storage{
		path:     path,
		capacity: capacity,
		cache:    cache.NewCache(capacity),
	}
}
