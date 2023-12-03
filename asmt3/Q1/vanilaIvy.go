package main

import (
	"flag"
	"fmt"
	performanceLogger "go_asmt/asmt1/asmt3"
	"sync"
	"time"
)

var start_time time.Time
var end_time time.Time

var timeout = 30 * time.Second

const (
	CM_ID         = -1 // id of central manager
	BACK_UP_CM_ID = -2
)

func main() {

	// log outputs for debugging purpose
	logger.Printf("===============START===============")

	// command line args
	num_servers := flag.Int("servers", 10, "number of servers")
	num_request_page := flag.Int("request_page", 3, "number of concurrent requests to make")

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
	for i := 0; i < *num_servers; i++ {
		for j := 0; j < *num_request_page; j++ {
			wg.Add(1)
			go func(i int) {
				if j%2 == 1 {
					servers[i].read(j)
				} else {
					servers[i].write(j)
				}
				wg.Done()
			}(i)
			time.Sleep(500 * time.Millisecond)
		}
	}

	// wait for all goroutines to finish
	wg.Wait()
	end_time = time.Now()
	elapsed_time := end_time.Sub(start_time) - timeout

	fmt.Println("Elapsed time: ", elapsed_time)
	logger.Println("Elapsed time: ", elapsed_time)
	fmt.Println("Refer to logs.txt for more logs.")

	performanceLogger.PerformanceLogger.Printf("| %d | %d | %d |\n", *num_servers, *num_request_page, elapsed_time)

	// var input string
	// // wait for the input, as otherwise, the program will not wait
	// fmt.Scanln(&input)
}
