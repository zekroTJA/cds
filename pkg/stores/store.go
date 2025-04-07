package stores

import (
	"io"
	"time"
)

type Metadata struct {
	Name         string
	IsDir        bool
	LastModified *time.Time
	Size         int64
	MimeType     string
}

type Store interface {
	Get(path string) (rc io.ReadCloser, metadata *Metadata, err error)
	List(path string) ([]*Metadata, error)
	IsNotExist(err error) bool
}

type StoreEntry struct {
	Entrypoint   string
	Listable     bool
	CacheControl string
	Store        Store
}
