package main

import (
	"flag"
	"fmt"
	performanceLogger "go_asmt/asmt1/asmt2"
	"sync"
	"time"
)

var start_time time.Time
var end_time time.Time
var time_mutex sync.Mutex
var timeout = 10 * time.Second

func main() {

	// log outputs for debugging purpose
	// log outputs for debugging purpose
	logger.Printf("===============START===============")

	// define the number of servers
	num_servers := flag.Int("servers", 10, "number of servers")

	// define the number of concurrent requests to make
	num_requests := flag.Int("requests", 10, "number of concurrent requests to make")
	flag.Parse()

	// set wg
	var wg sync.WaitGroup

	// initialize servers
	var servers = make([]*Server, *num_servers)
	for i := 0; i < *num_servers; i++ {
		servers[i] = NewServer(i, *num_servers, servers)
	}

	// start listening on all servers
	for i := 0; i < *num_servers; i++ {
		wg.Add(1)
		go func(i int) {
			servers[i].listen()
			wg.Done()
		}(i)
	}

	// simulate request
	start_time = time.Now()
	for i := 0; i < *num_requests; i++ {
		wg.Add(1)
		go func(i int) {
			servers[i].request()
			wg.Done()
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()
	elapsed_time := end_time.Sub(start_time)
	fmt.Println("Elapsed time: ", elapsed_time)
	logger.Println("Elapsed time: ", elapsed_time)

	performanceLogger.PerformanceLogger.Printf("| %d | %d |\n", *num_requests, elapsed_time)

	// var input string
	// // wait for the input, as otherwise, the program will not wait
	// fmt.Scanln(&input)
}
