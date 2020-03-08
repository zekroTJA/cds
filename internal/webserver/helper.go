package webserver

import (
	"encoding/json"
	"os"

	"github.com/valyala/fasthttp"
)

var (
	checksumSha1   = []byte("sha1")
	checksumSha256 = []byte("sha256")
	checksumMd5    = []byte("md5")

	headerAllow         = []byte("Allow")
	headerXForwardedFor = []byte("X-Forwarded-For")
	headerAllowValue    = []byte("GET,OPTIONS")

	contentTypeJSON      = []byte("application/json")
	contentTypeTextPlain = []byte("text/plain; charset=utf-8")
)

var errorMessages = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "not found",
	405: "method not allowed",
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func respondJSON(ctx *fasthttp.RequestCtx, statusCode int, v interface{}) (err error) {
	var data []byte

	if v != nil {
		data, err = json.Marshal(v)
		if err != nil {
			return err
		}
	}

	ctx.Response.Header.SetContentTypeBytes(contentTypeJSON)
	ctx.SetStatusCode(statusCode)
	_, err = ctx.Write(data)

	return
}

func (ws *WebServer) respondError(ctx *fasthttp.RequestCtx, statusCode int, errorMessage string) error {
	ctx.SetStatusCode(statusCode)

	pageFile, ok := ws.config.StatusPages[statusCode]
	if ok {
		if _, err := os.Stat(pageFile); err == nil {
			ctx.SendFile(pageFile)
			return nil
		}
	}

	if errorMessage == "" {
		errorMessage, _ = errorMessages[statusCode]
	}

	return respondJSON(ctx, statusCode, &ErrorResponse{
		Code:    statusCode,
		Message: errorMessage,
	})
}

func getIPAddr(ctx *fasthttp.RequestCtx) string {
	forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	if forwardedfor != nil && len(forwardedfor) > 0 {
		return string(forwardedfor)
	}

	return ctx.RemoteIP().String()
}
