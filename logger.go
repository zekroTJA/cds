package main

import (
	"log"
	"fmt"
)

const (
	typeInfo  = "INFO "
	typeWarn  = "WARN "
	typeError = "ERROR"
	typeFatal = "FATAL"
)

func printlog(logType string, stuff []interface{}) {
	message := " | " + logType + " | "
	for _, e := range stuff {
		message += fmt.Sprintf("%+v ", e)
	}
	if logType == typeFatal {
		log.Fatal(message)
	}
	log.Println(message)
}

func LogInfo(stuff ...interface{}) {
	printlog(typeInfo, stuff)
}

func LogWarn(stuff ...interface{}) {
	printlog(typeWarn, stuff)
}

func LogError(stuff ...interface{}) {
	printlog(typeError, stuff)
}

func LogFatal(stuff ...interface{}) {
	printlog(typeFatal, stuff)
}