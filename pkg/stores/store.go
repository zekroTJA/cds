package stores

import (
	"io"
	"time"
)

type PathEntry struct {
	Name  string
	IsDir bool
}

type Store interface {
	Get(path string) (rc io.ReadCloser, mimeType string, lastModified time.Time, err error)
	List(path string) ([]PathEntry, error)
	IsNotExist(err error) bool
}

type StoresEntry struct {
	Entrypoint   string
	Listable     bool
	CacheControl string
	Store        Store
}

type Stores []StoresEntry
