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

var timeout = 60 * time.Second
var heartbeat_interval = 1 * time.Second
var heartbeat_timeout = 2 * time.Second
var resend_timeout = 2 * time.Second

var completed_req int = 0
var mu sync.Mutex

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

	// command line args
	num_server := flag.Int("servers", 5, "number of servers")
	num_request := flag.Int("requests", 3, "number of requests for each server to make")
	num_faults_primary := flag.Int("faults", 0, "number of times that the primary fails")
	rejoin_primary := flag.Bool("rejoin", false, "whether primary rejoins after fails")
	fail_backup := flag.Bool("fail_backup", false, "whether the backup manager fails as well as the primary")

	flag.Parse()

	// set wg
	var wg sync.WaitGroup

	// initially, every server possess one page, the page number&content of it is its own id
	var managerRecords = make([]ManagerRecord, 0)
	for i := 0; i < *num_server; i++ {
		managerRecords = append(managerRecords, NewManagerRecord(i, make([]int, 0), i))
	}

	// initialize managers
	var servers = make([]*Server, *num_server)
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
	for i := 0; i < *num_server; i++ {
		serverRecords := make([]ServerRecord, 1)
		serverRecords = append(serverRecords, NewServerRecord(i, RW, NewPage(i, i)))
		servers[i] = NewServer(i, servers, primaryManager, managers, serverRecords)
	}

	// start listening on all servers
	for i := 0; i < *num_server; i++ {
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
		for i := 0; i < *num_faults_primary; i++ {
			// simulate primary down
			primaryManager.down()

			if *rejoin_primary {
				// simulate primary rejoin
				time.Sleep(3 * time.Second)
				primaryManager.rejoin()
			}

			if *fail_backup {
				// simulate backup down
				backupManager.down()

				// simulate backup rejoin
				time.Sleep(3 * time.Second)
				backupManager.rejoin()
			}
			// rest a while, don't fail so frequently
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// simulate request
	start_time = time.Now()
	for i := 0; i < *num_server; i++ {
		for j := 0; j < *num_request; j++ {
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
	completion_rate := float64(completed_req) / float64((*num_server)*(*num_request)) * 100

	fmt.Println("Elapsed time: ", elapsed_time)
	logger.Println("Elapsed time: ", elapsed_time)
	fmt.Printf("Request completion rate: %.0f%%\n", completion_rate)
	logger.Printf("Request completion rate: %.0f%%\n", completion_rate)
	fmt.Println("Refer to logs.txt for more logs.")

	performanceLogger.PerformanceLogger.Printf("| %d | %d | %d | %.0f%% |\n", *num_server, *num_request, elapsed_time, completion_rate)

	// var input string
	// // wait for the input, as otherwise, the program will not wait
	// fmt.Scanln(&input)
}
