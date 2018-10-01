package main

import (
	"os"
	"fmt"
	"flag"
)

var (
	configFile = flag.String("c", "config.yaml", "config file location")
	version    = flag.Bool("v", false, "show build version")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("cds v.%s\n© 2018 zekro Development\n", VERSION)
		return
	}

	LogInfo(fmt.Sprintf("Initializing CDS (v.%s)", VERSION))

	config, err := OpenConfig(*configFile)
	if os.IsNotExist(err) {
		err = CreateConfig(*configFile)
		if err != nil {
			LogFatal("Could not find config and failed creating config in current run directory:", err)
		}
		LogFatal("Could not find config.yaml. File was created in current run directory. Please configure the file and restart.")
	}
	LogInfo("Config loaded")
	
	mysql, err := NewMySql(
		config.MySql.Address, 
		config.MySql.Username, 
		config.MySql.Password, 
		config.MySql.Database,
		"dbscheme.sql")
	if err != nil {
		LogError("Could not connect to MySql database:", err)
	} else {
		LogInfo("Connected to MySql database")
	}
	defer mysql.Close()
	
	err = OpenWebServer(config.Port, mysql, config.DataPaths, config.TLS, config.Logging)
	if err != nil {
		LogFatal("Could not open WebServer:", err)
	}
}