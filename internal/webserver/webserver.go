package webserver

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"hash"
	"io"
	"os"

	"github.com/zekroTJA/cds/internal/logger"
	"github.com/zekroTJA/cds/internal/static"
	"github.com/zekroTJA/cds/internal/util"

	"github.com/valyala/fasthttp"
	"github.com/zekroTJA/cds/internal/config"
	"github.com/zekroTJA/cds/internal/database"
)

type WebServer struct {
	config *config.WebServer
	db     *database.MySQL

	server *fasthttp.Server
}

type fileCheck struct {
	status   int
	errMsg   string
	path     string
	fileName string
}

func NewWebServer(config *config.WebServer, db *database.MySQL) *WebServer {
	ws := &WebServer{
		config: config,
		db:     db,
	}

	ws.server = &fasthttp.Server{
		Handler: ws.handleRequest,
	}

	return ws
}

func (ws *WebServer) handleRequest(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "GET":
	case "OPTIONS":
		ctx.Response.Header.SetBytesKV(headerAllow, headerAllowValue)
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	default:
		ctx.Response.Header.SetBytesKV(headerAllow, headerAllowValue)
		ws.respondError(ctx, fasthttp.StatusMethodNotAllowed, "")
		return
	}

	if len(ctx.Path()) == 1 {
		if ws.config.RootRedirect != "" {
			ctx.Redirect(ws.config.RootRedirect, fasthttp.StatusSeeOther)
			return
		}
		respondJSON(ctx, fasthttp.StatusOK, map[string]string{
			"info":       "cds 2.0",
			"version":    static.AppVersion,
			"copyright":  "© 2019-2020 Ringo Hoffmann (zekro Development) [MAY NOT THE SERVER HOST]",
			"repository": "https://github.com/zekroTJA/cds",
			"licence":    "MIT",
		})
		return
	}

	ws.handleServeFile(ctx)

}

func (ws *WebServer) handleServeFile(ctx *fasthttp.RequestCtx) {
	fc := &fileCheck{
		status: fasthttp.StatusNotFound,
		path:   "/",
	}

	defer func() {
		if err := ws.db.RecordAccess(fc.path, fc.fileName, getIPAddr(ctx),
			string(ctx.UserAgent()), string(ctx.URI().FullURI()), ctx.Response.StatusCode()); err != nil {
			logger.Error("DATABASE :: failed recored access: %s", err.Error())
		}
	}()

	reqPath := ctx.Path()

	for _, storage := range ws.config.Storages {
		path := util.ConcatToString([]byte(storage), reqPath)
		fc = ws.checkFile(path)
		if fc.status < 400 {
			if ws.handleChecksums(ctx, fc) {
				return
			}
			ctx.SendFile(fc.path)
			return
		}
	}

	ws.respondError(ctx, fc.status, fc.errMsg)
}

func (ws *WebServer) checkFile(path string) *fileCheck {
	var status int
	var errMsg string

	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			status = fasthttp.StatusNotFound
		} else if os.IsPermission(err) {
			status = fasthttp.StatusForbidden
		} else {
			status = fasthttp.StatusInternalServerError
			errMsg = err.Error()
		}

		if path[0:2] == "./" {
			path = path[2:]
		}

		return &fileCheck{
			errMsg:   errMsg,
			status:   status,
			path:     path,
			fileName: "",
		}
	}

	if fi.IsDir() {
		return ws.checkFile(path + "/index.html")
	}

	return &fileCheck{
		errMsg:   errMsg,
		status:   fasthttp.StatusOK,
		path:     path,
		fileName: fi.Name(),
	}
}

func (ws *WebServer) handleChecksums(ctx *fasthttp.RequestCtx, fc *fileCheck) bool {
	var hasher hash.Hash

	checksum := ctx.QueryArgs().Peek("checksum")

	if checksum == nil {
		return false
	}

	if bytes.Equal(checksum, checksumMd5) {
		hasher = md5.New()
	} else if bytes.Equal(checksum, checksumSha1) {
		hasher = sha1.New()
	} else if bytes.Equal(checksum, checksumSha256) {
		hasher = sha256.New()
	} else {
		ws.respondError(ctx, fasthttp.StatusBadRequest,
			"unsupported hashing method")
		return true
	}

	if hasher == nil {
		return false
	}

	f, err := os.Open(fc.path)
	if err != nil {
		ws.respondError(ctx, fasthttp.StatusInternalServerError, err.Error())
		return true
	}

	defer f.Close()

	if _, err = io.Copy(hasher, f); err != nil {
		ws.respondError(ctx, fasthttp.StatusInternalServerError, err.Error())
		return true
	}

	hash := hasher.Sum(nil)

	ctx.Response.Header.SetContentTypeBytes(contentTypeTextPlain)
	ctx.SetStatusCode(fasthttp.StatusOK)
	hex.NewEncoder(ctx).Write(hash)
	// ctx.Write(hash)

	return true
}

func (ws *WebServer) ListenAndServeBlocking() error {
	tls := ws.config.TLS

	if tls.Enable {
		if tls.CertFile == "" || tls.KeyFile == "" {
			return errors.New("cert file and key file must be specified")
		}
		return ws.server.ListenAndServeTLS(ws.config.Addr, tls.CertFile, tls.KeyFile)
	}

	return ws.server.ListenAndServe(ws.config.Addr)
}
