package stores

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

type Local struct {
	basePath string
}

var _ Store = (*Local)(nil)

func NewLocal(basePath string) (*Local, error) {
	t := &Local{
		basePath: basePath,
	}

	return t, nil
}

func (t *Local) Get(path string) (io.ReadCloser, *Metadata, error) {
	path = t.path(path)

	stat, err := os.Stat(path)
	if err != nil {
		return nil, nil, err
	}

	if stat.IsDir() {
		return nil, nil, ErrNotExist
	}

	var metadata Metadata

	metadata.Name = stat.Name()
	metadata.Size = stat.Size()

	lastMod := stat.ModTime()
	metadata.LastModified = &lastMod

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	ext := filepath.Ext(path)
	metadata.MimeType = mime.TypeByExtension(ext)

	if metadata.MimeType == "" {
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		metadata.MimeType = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart)
		if err != nil {
			return nil, nil, err
		}
	}

	return f, &metadata, nil
}

func (t *Local) List(path string) ([]*Metadata, error) {
	entries, err := os.ReadDir(t.path(path))
	if err != nil {
		return nil, err
	}

	res := make([]*Metadata, 0, len(entries))
	for _, entry := range entries {
		var metadata Metadata
		metadata.Name = entry.Name()
		metadata.IsDir = entry.IsDir()

		info, err := entry.Info()
		if err == nil {
			lastMod := info.ModTime()
			metadata.LastModified = &lastMod
			metadata.Size = info.Size()
		}

		res = append(res, &metadata)
	}

	return res, nil
}

func (t *Local) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (t *Local) path(path string) string {
	return filepath.Join(t.basePath, path)
}
