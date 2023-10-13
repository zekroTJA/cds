package server

import (
	"embed"
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
	"html/template"
	"net/http"
)

//go:embed pages
var pages embed.FS

var tpl = template.Must(template.New("").ParseFS(pages, "pages/templates/*.html"))

func renderIndex(w http.ResponseWriter, dirName string, entries []stores.PathEntry) {
	err := tpl.ExecuteTemplate(w, "index.html", struct {
		DirName string
		Entries []stores.PathEntry
	}{
		DirName: dirName,
		Entries: entries,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed rendering index")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
