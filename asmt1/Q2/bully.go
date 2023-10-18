package main

// 1. Implement the above protocol of joint synchronization and election via GO (20 points)
// The coordinator initiates the synchronization process by sending message to all other machines.
// Upon receiving the message from the coordinator, each machine updates its local version of the data structure with the coordinator’s version.
// a new coordinator is chosen by the Bully algorithm. You can assume a fixed timeout to simulate the behaviour of detecting a fault.
// version working w/o bully. can send sync msg

import (
	"fmt"
	"time"
)

func (n *Node) replica_routine() {
	go n.rep_sync()
	go n.rep_elect()

	fmt.Printf("Node %d starts to execute replica_routine.\n", n.id)
	logger.Printf("Node %d starts to execute replica_routine.\n", n.id)
}

func (n *Node) coordinator_routine() {

	go n.coor_elect()
	go n.coor_sync()

	fmt.Printf("Node %d starts to execute coordinator_routine.\n", n.id)
	logger.Printf("Node %d starts to execute coordinator_routine.\n", n.id)
}

func (n *Node) execute() {
	for {
		if n.role == COORDINATOR {
			n.coordinator_routine()
		} else if n.role == REPLICA {
			n.replica_routine()
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

	// define the number of nodes
	const num_nodes int = 3
	const sync_interval = time.Second
	const timeout = 6 * time.Second // 2T(m) + T(p)

	// initialize (num_nodes-1) replicas
	var nodes = make([]*Node, num_nodes)
	for i := 0; i < num_nodes; i++ {
		nodes[i] = NewNode(i, sync_interval, timeout, REPLICA, nodes)
	}

	// make graph[highest] coordinator
	nodes[num_nodes-1].role = COORDINATOR

	// execute node go routines
	for i := 0; i < num_nodes; i++ {
		go nodes[i].execute()
	}

	// Simulate coordinator failure for once
	go func() {
		// var coorNode Node
		// for i := range nodes {
		// 	if .status == COORDINATOR {
		// 		coorNode = *nodes[i]
		// 		break
		// 	}
		// }
		time.Sleep(4 * time.Second)
		nodes[num_nodes-1].fail()
	}()

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
