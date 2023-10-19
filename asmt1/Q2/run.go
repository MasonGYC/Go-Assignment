package main

// 1. Implement the above protocol of joint synchronization and election via GO (20 points)
// The coordinator initiates the synchronization process by sending message to all other machines.
// Upon receiving the message from the coordinator, each machine updates its local version of the data structure with the coordinator’s version.
// a new coordinator is chosen by the Bully algorithm. You can assume a fixed timeout to simulate the behaviour of detecting a fault.
// version working w/o bully. can send sync msg

import (
	"flag"
	"fmt"
	"time"
)

// replica perform its tasks
func (n *Node) replica_tasks() {
	go n.rep_sync()
	go n.rep_elect()

	fmt.Printf("Node %d starts to execute replica_tasks.\n", n.id)
	logger.Printf("Node %d starts to execute replica_tasks.\n", n.id)
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
		if n.role == COORDINATOR {
			n.coordinator_tasks()
		} else if n.role == REPLICA {
			n.replica_tasks()
		}

		role_switch_msg := <-n.ch_role_switch
		fmt.Println(role_switch_msg)
		fmt.Printf("Node %d switch to role.\n", n.id)
		logger.Printf("Node %d switch role.\n", n.id)
	}

}

func main() {

	// log outputs for debugging purpose
	logger.Println("===============START===============")

	// define the number of nodes, pass by flag
	num_nodes := flag.Int("nodes", 3, "number of nodes")
	sync_interval_second := flag.Int("sync", 1, "the time interval to sync in second")
	timeout_second := flag.Int("timeout", 6, "2T(m) + T(p) in second")
	flag.Parse()

	var sync_interval = time.Duration(*sync_interval_second) * time.Second
	var timeout = time.Duration(*timeout_second) * time.Second

	// initialize (num_nodes-1) replicas
	var nodes = make([]*Node, *num_nodes)
	for i := 0; i < *num_nodes; i++ {
		nodes[i] = NewNode(i, sync_interval, timeout, REPLICA, nodes)
	}

	// make graph[highest] coordinator
	nodes[*num_nodes-1].role = COORDINATOR

	// execute node go routines
	for i := 0; i < *num_nodes; i++ {
		go nodes[i].execute()
	}

	// Simulate coordinator failure after 4 seconds
	go func() {
		time.Sleep(4 * time.Second)
		for i := range nodes {
			if nodes[i].role == COORDINATOR {
				nodes[i].fail_during_election()
				break
			}
		}
	}()

	// Simulate coordinator candicate (newly seleted coordinator) failure duing broadcasting
	go func() {
		time.Sleep(4 * time.Second)
		for {
			for i := range nodes {
				if nodes[i].state == BROADCATING {
					nodes[i].fail()
					return
				}
			}
		}

	}()

	// // Simulate non-coordinator failure duing someone else broadcasting
	// go func() {
	// 	for {
	// 		time.Sleep(10 * time.Second) //check every 5 seconds
	// 		for i := range nodes {
	// 			if nodes[i].state == BROADCATING {
	// 				for {
	// 					random_node_idx := rand.Intn(len(nodes))
	// 					if random_node_idx != i && nodes[random_node_idx].state != DOWN {
	// 						nodes[random_node_idx].fail()
	// 						return
	// 					}

	// 				}
	// 			}
	// 		}
	// 	}

	// }()

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
