package performanceLogger

import (
	"log"
	"os"
)

var PerformanceLogger *log.Logger

func init() {
	// log outputs for debugging purpose
	file, err := os.OpenFile("performance_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	PerformanceLogger = log.New(file, "", 0)
}
