package server

import (
	"embed"
	"html/template"
	"io"
	"net/http"
	"path"

	"github.com/dustin/go-humanize"
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
)

//go:embed pages
var pages embed.FS

var tpl = template.Must(template.New("").
	Funcs(template.FuncMap{
		"humanBytes": func(size int64) string { return humanize.Bytes(uint64(size)) },
	}).
	ParseFS(pages, "pages/templates/*.html"))

func renderIndex(w http.ResponseWriter, dirName string, entries []*stores.Metadata, sortBy string, ascending bool) {
	err := tpl.ExecuteTemplate(w, "index.html", struct {
		DirName   string
		Entries   []*stores.Metadata
		SortBy    string
		Ascending bool
	}{
		DirName:   dirName,
		Entries:   entries,
		SortBy:    sortBy,
		Ascending: ascending,
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
