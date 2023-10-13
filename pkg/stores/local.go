package stores

import (
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

func (t *Local) Get(path string) (rc io.ReadCloser, mimeType string, lastModified time.Time, err error) {
	path = t.path(path)

	stat, err := os.Stat(path)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	if stat.IsDir() {
		return nil, "", time.Time{}, fs.ErrNotExist
	}

	lastModified = stat.ModTime()

	f, err := os.Open(path)
	if err != nil {
		return nil, "", time.Time{}, err
	}

	ext := filepath.Ext(path)
	mimeType = mime.TypeByExtension(ext)

	if mimeType == "" {
		var buf [512]byte
		n, _ := io.ReadFull(f, buf[:])
		mimeType = http.DetectContentType(buf[:n])
		_, err := f.Seek(0, io.SeekStart) // rewind to output whole file
		if err != nil {
			return nil, "", time.Time{}, err
		}
	}

	return f, mimeType, lastModified, nil
}

func (t *Local) List(path string) ([]PathEntry, error) {
	entries, err := os.ReadDir(t.path(path))
	if err != nil {
		return nil, err
	}

	res := make([]PathEntry, 0, len(entries))
	for _, e := range entries {
		res = append(res, PathEntry{
			Name:  e.Name(),
			IsDir: e.IsDir(),
		})
	}

	return res, nil
}

func (t *Local) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (t *Local) path(path string) string {
	return filepath.Join(t.basePath, path)
}
