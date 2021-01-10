package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/zekroTJA/cds/internal/util"
)

type Main struct {
	WebServer *WebServer `json:"web_server"`
	MySQL     *MySQL     `json:"mysql"`
}

type WebServer struct {
	Addr         string           `json:"addr"`
	RootRedirect string           `json:"root_redirect"`
	TLS          *WebServerTSL    `json:"tls"`
	Upload       *WebServerUpload `json:"upload"`
	StatusPages  map[int]string   `json:"status_pages"`
	Storages     []string         `json:"storages"`
}

type WebServerTSL struct {
	Enable   bool   `json:"enable"`
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type WebServerUpload struct {
	Enable         bool   `json:"enable"`
	Secret         string `json:"secret"`
	AllowOverwrite bool   `json:"allow_overwrite"`
	Storage        string `json:"storage"`
	MaxSizeBytes   int    `json:"max_size_bytes"`
}

type MySQL struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

func Open(loc string) (*Main, error) {
	data, err := ioutil.ReadFile(loc)
	if os.IsNotExist(err) {
		err = createDefault(loc)
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	cfg := new(Main)
	err = yaml.Unmarshal(data, cfg)
	return cfg, err
}

func createDefault(loc string) error {
	randSecret, err := util.GetRandBase64Str(32)
	if err != nil {
		return err
	}

	def := &Main{
		WebServer: &WebServer{
			Addr:         ":80",
			RootRedirect: "",
			StatusPages: map[int]string{
				404: "./web/pages/404.html",
				401: "./web/pages/401.html",
			},
			TLS: &WebServerTSL{
				Enable: false,
			},
			Upload: &WebServerUpload{
				Enable:         false,
				Secret:         randSecret,
				AllowOverwrite: false,
				Storage:        "",
				MaxSizeBytes:   2 * 1024 * 1024 * 1024,
			},
		},
		MySQL: &MySQL{
			Address:  "localhost",
			Database: "cds",
			Username: "cds",
		},
	}

	data, err := yaml.Marshal(def)

	basePath := path.Dir(loc)
	if _, err = os.Stat(basePath); os.IsNotExist(err) {
		err = os.MkdirAll(basePath, 0750)
		if err != nil {
			return err
		}
	}
	err = ioutil.WriteFile(loc, data, 0750)
	return err
}
