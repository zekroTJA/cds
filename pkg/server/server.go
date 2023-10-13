package server

import (
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/log"
	"io"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	stores stores.Stores
}

func New(storeMap stores.Stores) (t *Server, err error) {
	t = &Server{}

	t.stores = storeMap

	return t, nil
}

func (t *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, e := range t.stores {
		if strings.HasPrefix(r.URL.Path, e.Entrypoint) {
			sub := subPath(r.URL.Path, e.Entrypoint)
			if sub == "" {
				continue
			}

			rc, mimeType, lastModified, err := e.Store.Get(sub)
			if err != nil {
				if e.Store.IsNotExist(err) {
					continue
				}
				handleError(err, w, r, "failed serving object")
				return
			}
			defer rc.Close()

			cacheControl := e.CacheControl
			if cacheControl == "" {
				cacheControl = "public, max-age=2592000, must-revalidate"
			}

			w.Header().Set("Cache-Control", cacheControl)
			w.Header().Set("Content-Type", mimeType)
			w.Header().Set("Last-Modified", lastModified.Format(time.RFC1123))

			_, err = io.Copy(w, rc)
			if err != nil {
				log.Warn().Err(err).Fields("path", r.URL.Path).Msg("failed writing response")
			}

			return
		}
	}

	var entries []stores.PathEntry
	for _, st := range t.stores {
		if !st.Listable || !strings.HasPrefix(r.URL.Path, st.Entrypoint) {
			continue
		}

		sub := subPath(r.URL.Path, st.Entrypoint)
		e, err := st.Store.List(sub)
		if err != nil {
			if !st.Store.IsNotExist(err) {
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

func subPath(path, entrypoint string) string {
	sub := path[len(entrypoint):]
	if sub != "" && sub[0] == '/' {
		sub = sub[1:]
	}
	return sub
}
