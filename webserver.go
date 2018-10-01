package main

import (
	"os"
	"time"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
)

var mysql *MySql
var filePaths []string
var confLogging *ConfigLogging

type WebServerCert struct {
	CertFile string `yaml:"certfile"`
	KeyFile  string	`yaml:"keyfile"`
}

type WebServerError struct {
	Code    int		`json:"code"`
	Message string  `json:"message"`
}

func writeError(w http.ResponseWriter, code int, message string) {
	err := &WebServerError{code, message}
	bdata, _ := json.MarshalIndent(err, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(bdata)
}

func fileServerHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	pathSplit := strings.Split(urlPath, "/")
	fileName := pathSplit[len(pathSplit)-1]
	status := 304

	defer func() {
		if mysql != nil && confLogging.RequestLog {
			mysql.Query(
				"INSERT INTO requestLog (address, userAgent, url, code) VALUES (?, ?, ?, ?)",
				r.RemoteAddr, r.Header.Get("User-Agent"), r.URL.String(), status)
		}
	}()

	for _, path := range filePaths {
		stat, err := os.Stat(path + urlPath)
		if os.IsPermission(err) || (stat != nil && stat.IsDir()) {
			writeError(w, 403, "permission denied")
			status = 403
			return
		}
		if !os.IsNotExist(err) {
			http.ServeFile(w, r, path + urlPath)
			if mysql != nil && confLogging.AccessCounts {
				rows, err := mysql.Query("SELECT * FROM accessStats WHERE fullPath = ?;", urlPath)
				if err == nil {
					var cnt int
					for rows.Next() {
						cnt++
					}
					if cnt > 0 {
						mysql.Query(
							"UPDATE accessStats SET accesses = accesses+1, lastAccess = ? WHERE fullPath = ?;",
							time.Now(), urlPath)
					} else {
						mysql.Query(
							"INSERT INTO accessStats (fullPath, fileName, accesses) VALUES (?, ?, 1);",
							urlPath, fileName)
					}
				}
			}
			return
		}
	}

	status = 404
	writeError(w, 404, "not found")
}

func OpenWebServer(port string, db *MySql, dataFilePaths []string, cert *WebServerCert, logging *ConfigLogging) error {
	certEnabled := cert.KeyFile != "" && cert.CertFile != ""

	router := mux.NewRouter()

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		return
	})

	router.Methods("GET").HandlerFunc(fileServerHandler)

	http.Handle("/", router)

	mysql = db
	filePaths = dataFilePaths
	confLogging = logging

	if certEnabled {
		LogInfo("WebServer listening in TLS mode on port", port)
		return http.ListenAndServeTLS(":" + port, cert.CertFile, cert.KeyFile, nil)
	} else {
		LogWarn("WebServer listening in NON TLS mode on port", port)
		return http.ListenAndServe(":" + port, nil)
	}
}