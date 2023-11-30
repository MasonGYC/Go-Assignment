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

const (
	CM_ID         = -1 // id of central manager
	BACK_UP_CM_ID = -2
)

func main() {

	// log outputs for debugging purpose
	logger.Printf("===============START===============")

	// define the number of servers
	num_servers := flag.Int("servers", 4, "number of servers")

	// define the number of concurrent requests to make
	num_requests := flag.Int("requests", 2, "number of concurrent requests to make")
	flag.Parse()

	// set wg
	var wg sync.WaitGroup

	// initially, every server possess one page, the page number&content of it is its own id
	var managerRecords = make([]ManagerRecord, 0)
	for i := 0; i < *num_servers; i++ {
		managerRecords = append(managerRecords, NewManagerRecord(i, make([]int, 0), i))
	}

	// initialize servers and manager
	var servers = make([]*Server, *num_servers)
	var manager = NewManager(servers, managerRecords)
	for i := 0; i < *num_servers; i++ {
		serverRecords := make([]ServerRecord, 1)
		serverRecords = append(serverRecords, NewServerRecord(i, RW, NewPage(i, i)))
		servers[i] = NewServer(i, servers, manager, serverRecords)
	}

	// start listening on all servers
	for i := 0; i < *num_servers; i++ {
		wg.Add(1)
		go func(i int) {
			servers[i].listen()
			wg.Done()
		}(i)
	}

	wg.Add(1)
	go func() {
		manager.start()
		wg.Done()
	}()

	// simulate request
	start_time = time.Now()
	for i := 0; i < *num_requests; i++ {
		wg.Add(1)
		go func(i int) {
			servers[i].write(3)
			wg.Done()
		}(i)
	}

	// wait for all goroutines to finish
	wg.Wait()
	end_time = time.Now()
	elapsed_time := end_time.Sub(start_time)

	fmt.Println("Elapsed time: ", elapsed_time)
	logger.Println("Elapsed time: ", elapsed_time)
	fmt.Println("Refer to logs.txt for more logs.")

	performanceLogger.PerformanceLogger.Printf("| %d | %d |\n", *num_requests, elapsed_time)

	// var input string
	// // wait for the input, as otherwise, the program will not wait
	// fmt.Scanln(&input)
}
