package main

import (
	"os"
	"fmt"
)

func main() {

	LogInfo(fmt.Sprintf("Initializing CDS (v.%s)", VERSION))

	config, err := OpenConfig("config.yaml")
	if os.IsNotExist(err) {
		err = CreateConfig("config.yaml")
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