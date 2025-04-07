package server

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
)

//go:embed pages
var pages embed.FS

var tpl = template.Must(template.New("").ParseFS(pages, "pages/templates/*.html"))

func renderIndex(w http.ResponseWriter, dirName string, entries []*stores.Metadata) {
	err := tpl.ExecuteTemplate(w, "index.html", struct {
		DirName string
		Entries []*stores.Metadata
	}{
		DirName: dirName,
		Entries: entries,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed rendering index")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func servePage(w http.ResponseWriter, status int, pth string) {
	p, err := pages.Open(path.Join("pages", pth))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = io.Copy(w, p)
}
