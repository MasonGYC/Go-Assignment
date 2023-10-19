package main

import (
	"log"
	"os"
)

var logger *log.Logger

func init() {
	// log outputs for debugging purpose
	file, err := os.OpenFile("logs1.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
}
