package storage

type Cacher interface {
	Init() error
	Upload(name string, data []byte) error
	Download(name string) ([]byte, error)
}
