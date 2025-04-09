package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zekroTJA/cds/pkg/router"
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
)

type ResponseType string

const (
	ResponseTypeJSON ResponseType = "json"
	ResponseTypeHTML ResponseType = "html"
)

type ErrorModel struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

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
		handleNotFound(r, w)
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

		handleIndex(r, w, entries)
		return
	}

	handleNotFound(r, w)
}

func (t *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, t)
}

func handleIndex(r *http.Request, w http.ResponseWriter, entries []*stores.Metadata) {
	switch responseType(r) {
	case ResponseTypeHTML:
		renderIndex(w, r.URL.Path, entries)
	case ResponseTypeJSON:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(entries)
	}
}

func handleError(err error, w http.ResponseWriter, r *http.Request, msg string) {
	log.Error().Err(err).Fields("path", r.URL.Path).Msg(msg)
	switch responseType(r) {
	case ResponseTypeHTML:
		servePage(w, http.StatusInternalServerError, "500.html")
	case ResponseTypeJSON:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(ErrorModel{Status: http.StatusInternalServerError, Message: "An unexpected internal server error has occurred"})
	}

}

func handleNotFound(r *http.Request, w http.ResponseWriter) {
	switch responseType(r) {
	case ResponseTypeHTML:
		servePage(w, http.StatusNotFound, "404.html")
	case ResponseTypeJSON:
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(ErrorModel{Status: http.StatusNotFound, Message: "The requested resource is not or no more available"})
	}
}

func responseType(r *http.Request) ResponseType {
	typ := ResponseType(r.URL.Query().Get("format"))
	if typ == ResponseTypeJSON {
		return ResponseTypeJSON
	}
	if typ == ResponseTypeHTML {
		return ResponseTypeHTML
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "text/html") {
		return ResponseTypeHTML
	}

	return ResponseTypeJSON
}
