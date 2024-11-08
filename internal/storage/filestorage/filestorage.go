package filestorage

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/randomurban/image-previewer/internal/storage"
	"github.com/randomurban/image-previewer/internal/storage/cache"
)

var _ storage.Cacher = (*FileCache)(nil)

type FileCache struct {
	cache    cache.Cache
	capacity int
	path     string
}

func NewStorage(path string, capacity int) storage.Cacher {
	return &FileCache{
		path:     path,
		capacity: capacity,
		cache:    cache.NewCache(capacity),
	}
}

func (s *FileCache) Init() error {
	_, err := os.Stat(s.path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(s.path, 0o755)
		if err != nil {
			return fmt.Errorf("creating storage directory: %w", err)
		}
	}

	dir, err := os.ReadDir(s.path)
	if err != nil {
		return fmt.Errorf("reading storage directory: %w", err)
	}

	files := make([]string, 0, len(dir))
	for _, file := range dir {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}

	if len(files) != 0 {
		if len(files) > s.capacity {
			log.Printf("too many files (%d) in storage (capacity=%d)", len(files), s.capacity)
		}
		// добавить все файлы в cache
		for _, filename := range files {
			err := s.setAndDelete(filename)
			if err != nil {
				return fmt.Errorf("init cache: %w", err)
			}
		}
	}
	return nil
}

func (s *FileCache) Upload(name string, data []byte) error {
	err2 := s.setAndDelete(name)
	if err2 != nil {
		return err2
	}

	file, err := os.Create(path.Join(s.path, name))
	if err != nil {
		return fmt.Errorf("creating file in storage: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("closing file in storage: %v", err)
		}
	}(file)

	_, err = io.Copy(file, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("writing file in storage: %w", err)
	}
	return nil
}

func (s *FileCache) setAndDelete(name string) error {
	deleted, wasInCache := s.cache.Set(cache.Key(name), name)
	if wasInCache {
		log.Printf("filename in cache: %s", name)
	}
	if deleted != nil {
		deletedName, ok := deleted.(string)
		if ok {
			err := s.deleteFile(deletedName)
			if err != nil {
				return fmt.Errorf("deleting file in storage: %w", err)
			}
		}
	}
	return nil
}

func (s *FileCache) deleteFile(deleted string) error {
	filename := path.Join(s.path, deleted)
	err := os.Remove(filename)
	if err != nil {
		return fmt.Errorf("remove file: %w", err)
	}
	log.Printf("delete file in storage: %s", filename)
	return nil
}

func (s *FileCache) Download(name string) ([]byte, error) {
	val, wasInCache := s.cache.Get(cache.Key(name))
	if !wasInCache {
		log.Printf("filename not in cache: %s", name)
		return nil, nil
	}
	nameFromCache, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("invalid data in cache")
	}
	if name != nameFromCache {
		log.Printf("filename in cache by %s: %s", name, nameFromCache)
	}
	data, err := os.ReadFile(path.Join(s.path, nameFromCache))
	if err != nil {
		return nil, fmt.Errorf("reading file from storage: %w", err)
	}
	return data, nil
}
