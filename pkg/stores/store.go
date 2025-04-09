package stores

import (
	"io"
	"time"
)

type Metadata struct {
	Name         string     `json:"name"`
	IsDir        bool       `json:"is_dir,omitempty"`
	LastModified *time.Time `json:"last_modified,omitempty"`
	Size         int64      `json:"size,omitempty"`
	MimeType     string     `json:"mime_type,omitempty"`
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
