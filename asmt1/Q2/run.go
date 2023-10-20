package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

// worker perform its tasks
func (n *Node) worker_tasks() {
	go n.worker_sync()
	go n.worker_elect()

	fmt.Printf("Node %d starts to execute worker_tasks.\n", n.id)
	logger.Printf("Node %d starts to execute worker_tasks.\n", n.id)
}

// coordinator perform its tasks
func (n *Node) coordinator_tasks() {

	go n.coor_elect()
	go n.coor_sync()

	fmt.Printf("Node %d starts to execute coordinator_tasks.\n", n.id)
	logger.Printf("Node %d starts to execute coordinator_tasks.\n", n.id)
}

// node execute different tasks according to its role.
// switch role if get notification from ch_role_switch channel.
func (n *Node) execute() {
	for {

		fmt.Printf("execute: Node %d's state: %s.\n", n.id, n.state)
		logger.Printf("execute Node %d's state: %s.\n", n.id, n.state)

		if n.state != DOWN {
			if n.role == COORDINATOR {
				n.coordinator_tasks()
			} else if n.role == WORKER {
				n.worker_tasks()
			}
		}

		<-n.ch_role_switch
		fmt.Printf("Node %d switch to role.\n", n.id)
		logger.Printf("Node %d switch role.\n", n.id)
	}

}

func main() {

	// log outputs for debugging purpose
	logger.Println("===============START===============")

	// define cmd flags
	num_nodes := flag.Int("nodes", 3, "number of nodes")
	sync_interval_second := flag.Int("sync", 1, "the time interval to sync in second")
	timeout_second := flag.Int("timeout", 6, "2T(m) + T(p) in second")

	fail_coordinator := flag.Bool("failcoor", false, "if set to true, simulate coordinator fails")
	fail_worker := flag.Bool("failworker", false, "if set to true, simulate coordinator fails")
	fail_coor_during_broadcasting := flag.Bool("failcoorvic", false, "if set to true, simulate newly elected coordinator fails while announcing")
	fail_worker_during_broadcasting := flag.Bool("failworkervic", false, "if set to true, simulate worker node fails while announcing")

	flag.Parse()

	var sync_interval = time.Duration(*sync_interval_second) * time.Second
	var timeout = time.Duration(*timeout_second) * time.Second

	// initialize (num_nodes-1) workers
	var nodes = make([]*Node, *num_nodes)
	for i := 0; i < *num_nodes; i++ {
		nodes[i] = NewNode(i, sync_interval, timeout, WORKER, nodes)
	}

	// make graph[highest] coordinator
	nodes[*num_nodes-1].role = COORDINATOR

	// execute node go routines
	for i := 0; i < *num_nodes; i++ {
		go nodes[i].execute()
	}

	// Simulate coordinator failure after 4 seconds
	if *fail_coordinator {
		go func() {
			time.Sleep(4 * time.Second)
			for i := range nodes {
				if nodes[i].role == COORDINATOR {
					nodes[i].fail()
					break
				}
			}
		}()
	}

	// Simulate worker failure after 4 seconds
	if *fail_worker {
		go func() {
			time.Sleep(4 * time.Second)
			for {
				random_node_idx := rand.Intn(len(nodes))
				if nodes[random_node_idx].role == WORKER {
					nodes[random_node_idx].fail()
					return
				}
			}
		}()
	}

	// Simulate coordinator candicate (newly seleted coordinator) failure duing broadcasting
	if *fail_coor_during_broadcasting {
		go func() {
			time.Sleep(4 * time.Second)
			for {
				for i := range nodes {
					if nodes[i].state == BROADCATING {
						nodes[i].coordinator_fail_during_broadcasting()
						return
					}
				}
			}
		}()
	}

	// Simulate non-coordinator failure duing someone else broadcasting
	if *fail_worker_during_broadcasting {
		go func() {
			time.Sleep(4 * time.Second)
			for {
				for i := range nodes {
					if nodes[i].state == BROADCATING {
						for {
							random_node_idx := rand.Intn(len(nodes))
							if random_node_idx != i && nodes[random_node_idx].state != DOWN {
								nodes[random_node_idx].worker_fail_during_broadcasting()
								return
							}
						}
					}
				}
			}
		}()
	}

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
