package main

import (
	"flag"
	"fmt"
	performanceLogger "go_asmt/asmt1/asmt3"
	"math/rand"
	"sync"
	"time"
)

var start_time time.Time
var end_time time.Time

var timeout = 5 * time.Second
var heartbeat_interval = 1 * time.Second
var heartbeat_timeout = 2 * time.Second
var resend_timeout = 2 * time.Second

// roles
const (
	CM_ID         int    = -1 // id of central manager
	BACK_UP_CM_ID int    = -2
	PRIMARY       string = "PRIMARY"
	BACKUP        string = "BACKUP"
)

func main() {

	// log outputs for debugging purpose
	logger.Printf("===============START===============")

	// define the number of servers
	num_servers := flag.Int("servers", 10, "number of servers")

	// define the number of concurrent requests to make
	num_requests := flag.Int("requests", 8, "number of concurrent requests to make")

	// define the number of tiems that the primary fails
	// primary_fail_times := flag.Int("requests", 8, "number of concurrent requests to make")
	flag.Parse()

	// set wg
	var wg sync.WaitGroup

	// initially, every server possess one page, the page number&content of it is its own id
	var managerRecords = make([]ManagerRecord, 0)
	for i := 0; i < *num_servers; i++ {
		managerRecords = append(managerRecords, NewManagerRecord(i, make([]int, 0), i))
	}

	// initialize managers
	var servers = make([]*Server, *num_servers)
	var managers = make([]*Manager, 2)
	var primaryManager *Manager
	var backupManager *Manager

	primaryManager = NewManager(CM_ID, servers, managerRecords, PRIMARY)
	backupManager = NewManager(BACK_UP_CM_ID, servers, managerRecords, BACKUP)

	primaryManager.primaryManager = primaryManager
	primaryManager.backupManager = backupManager
	backupManager.primaryManager = primaryManager
	backupManager.backupManager = backupManager

	managers[0] = primaryManager
	managers[1] = backupManager

	// initialize servers
	for i := 0; i < *num_servers; i++ {
		serverRecords := make([]ServerRecord, 1)
		serverRecords = append(serverRecords, NewServerRecord(i, RW, NewPage(i, i)))
		servers[i] = NewServer(i, servers, primaryManager, managers, serverRecords)
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
		primaryManager.start()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		backupManager.start()
		wg.Done()
	}()

	// simulate faults of managers
	time.Sleep(time.Second)
	go func() {
		// // simulate primary down
		// time.Sleep(100 * time.Millisecond)
		// primaryManager.down()
		// // simulate primary rejoin
		// time.Sleep(2 * time.Second)
		// primaryManager.rejoin()
		time.Sleep(100 * time.Millisecond)
		backupManager.down()
		// simulate primary rejoin
		time.Sleep(2 * time.Second)
		backupManager.rejoin()
	}()

	// simulate request
	start_time = time.Now()
	for i := 0; i < *num_requests; i++ {
		wg.Add(1)
		go func(i int) {
			randReq := rand.Intn(2)
			randPage := rand.Intn(*num_servers)
			if randReq == 0 {
				servers[i].read(randPage)
			} else {
				servers[i].write(randPage)
			}
			wg.Done()
		}(i)
		time.Sleep(500 * time.Millisecond)
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
