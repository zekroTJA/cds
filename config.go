package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigLogging struct {
	AccessCounts bool `yaml:"accesscounts"`
	RequestLog   bool `yaml:"requestlog"`
}

type Config struct {
	Port		 string			`yaml:"port"`
	DataPaths    []string 		`yaml:"datapaths"`
	TLS          *WebServerCert	`yaml:"tsl"`
	MySql struct{
		Address  string 		`yaml:"address"`
		Username string 		`yaml:"username"`
		Password string 		`yaml:"password"`
		Database string 		`yaml:"database"`
	}							`yaml:"mysql"`
	Logging      *ConfigLogging `yaml:"logging"`
}

func OpenConfig(path string) (*Config, error) {
	config := new(Config)
	fhandler, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(fhandler)
	err = decoder.Decode(config)
	return config, err
}

func CreateConfig(path string) error {
	config := new(Config)
	config.Port = "443"
	config.TLS = new(WebServerCert)
	config.Logging = new(ConfigLogging)
	
	fhandler, err := os.Create(path)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(fhandler)
	defer encoder.Close()
	err = encoder.Encode(config)
	return err
}