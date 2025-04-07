package server

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/cds/pkg/router"
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
)

type Server struct {
	router *router.Router[stores.StoreEntry]
}

func New(storeEntries []stores.StoreEntry) (t *Server, err error) {
	t = &Server{}

	t.router = &router.Router[stores.StoreEntry]{}

	for _, entry := range storeEntries {
		t.router.Add(entry.Entrypoint, entry)
	}

	return t, nil
}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	storeEntries, sub, ok := t.router.Match(r.URL.Path)
	if !ok {
		handleNotFound(w)
		return
	}

	if sub != "" {
		for _, store := range storeEntries {
			rc, metadata, err := store.Store.Get(sub)
			if err != nil {
				if store.Store.IsNotExist(err) {
					continue
				}
				handleError(err, w, r, "failed serving object")
				return
			}
			defer rc.Close()

			cacheControl := store.CacheControl
			if cacheControl == "" {
				cacheControl = "public, max-age=2592000, must-revalidate"
			}
			w.Header().Set("Cache-Control", cacheControl)

			if metadata.MimeType != "" {
				w.Header().Set("Content-Type", metadata.MimeType)
			}

			if metadata.LastModified != nil {
				w.Header().Set("Last-Modified", metadata.LastModified.Format(time.RFC1123))
			}

			_, err = io.Copy(w, rc)
			if err != nil {
				log.Warn().Err(err).Fields("path", r.URL.Path).Msg("failed writing response")
			}

			return
		}
	}

	var entries []*stores.Metadata
	for _, store := range storeEntries {
		if !store.Listable {
			continue
		}

		e, err := store.Store.List(sub)
		if err != nil {
			if !store.Store.IsNotExist(err) {
				handleError(err, w, r, "failed listing storage entries")
			}
			continue
		}

		entries = append(entries, e...)
	}

	if len(entries) != 0 {
		if !strings.HasSuffix(r.URL.Path, "/") {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		renderIndex(w, r.URL.Path, entries)
		return
	}

	handleNotFound(w)
}

func (t *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, t)
}

func handleError(err error, w http.ResponseWriter, r *http.Request, msg string) {
	log.Error().Err(err).Fields("path", r.URL.Path).Msg(msg)
	servePage(w, http.StatusInternalServerError, "500.html")
}

func handleNotFound(w http.ResponseWriter) {
	servePage(w, http.StatusNotFound, "404.html")
}
