package main

import (
	"flag"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/zekroTJA/cds/pkg/config"
	"github.com/zekroTJA/cds/pkg/server"
	"github.com/zekroTJA/cds/pkg/stores"
	"github.com/zekrotja/rogu/level"
	"github.com/zekrotja/rogu/log"
)

var (
	flagConfig = flag.String("c", "config.toml", "config file location")
)

func main() {
	flag.Parse()

	cfg, err := config.Parse(*flagConfig, "CDS_", config.Default)
	if err != nil {
		log.Fatal().Err(err).Msg("failed parsing config")
	}

	lvl, ok := level.FromString(cfg.Logging.Level)
	if !ok {
		log.Fatal().Field("level", cfg.Logging.Level).Msg("invalid log level")
	}
	log.SetLevel(lvl)

	log.Debug().Msgf("Config: %s", spew.Sdump(cfg))

	storeList := make([]stores.StoreEntry, 0, len(cfg.Stores))
	for _, st := range cfg.Stores {
		var (
			store stores.Store
			err   error
		)

		typ := strings.ToLower(string(st.Type))

		switch config.StoreType(typ) {
		case config.StoreTypeLocal:
			store, err = stores.NewLocal(st.Path)
		case config.StoreTypeS3:
			store, err = stores.NewS3(
				st.Endpoint, st.AccessKey, st.SecretKey, st.Region, st.Bucket, st.Path, st.Secure)
		default:
			log.Error().Field("type", typ).Msg("invalid or unsupported store type")
			continue
		}

		if err != nil {
			log.Error().Err(err).Fields("type", typ, "entrypoint", st.Entrypoint).
				Msg("failed initializing storage")
			continue
		}

		storeList = append(storeList, stores.StoreEntry{
			Entrypoint:   prefixEntrypoint(st.Entrypoint),
			Listable:     st.Listable,
			CacheControl: st.CacheControl,
			Store:        store,
		})
	}

	if ln := len(storeList); ln == 0 {
		log.Warn().Msg("No stores have been initialized")
	} else {
		log.Info().Field("n", len(storeList)).Msg("Stores initialized")
	}

	srv, err := server.New(storeList)
	if err != nil {
		log.Fatal().Err(err).Msg("failed initializing web server")
	}

	log.Info().Field("address", cfg.Address).Msg("Starting listening and serving ...")
	err = srv.ListenAndServe(cfg.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed binding web server")
	}
}

func prefixEntrypoint(e string) string {
	if e == "" || e[0] != '/' {
		return "/" + e
	}
	return e
}
