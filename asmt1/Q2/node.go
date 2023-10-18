package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	id                   int
	ch_sync              chan Message
	ch_elect             chan Message
	ch_stop_elect        chan Message
	ch_role_switch       chan Message
	data                 Data
	sync_interval        time.Duration // heartbeat time interval, timeout
	timeout              time.Duration // if does not receive ack for this time, self-elect
	role                 string        // coor or rep
	nodes                []*Node
	coordinator_id       int
	coordinator_msg_time time.Time
	state                string
	mutex                sync.Mutex
}

func NewNode(id int, sync_interval time.Duration, timeout time.Duration, role string, nodes []*Node) *Node {
	return &Node{
		id:                   id,
		ch_sync:              make(chan Message),
		ch_elect:             make(chan Message),
		ch_stop_elect:        make(chan Message),
		ch_role_switch:       make(chan Message),
		data:                 Data{},
		sync_interval:        sync_interval,
		timeout:              timeout,
		role:                 role,
		nodes:                nodes,
		coordinator_id:       len(nodes) - 1, //highest
		coordinator_msg_time: time.Now(),
		state:                NORMAL,
		mutex:                sync.Mutex{},
	}
}

// send election message to higher id nodes
func (n *Node) self_elect() {

	fmt.Printf("Node %d started an election.\n", n.id)
	logger.Printf("Node %d started an election.\n", n.id)

	select {
	case msg := <-n.ch_stop_elect:
		fmt.Printf("Node %d %s.\n", n.id, msg.message)
		logger.Printf("Node %d %s.\n", n.id, msg.message)
		return
	default:
		nodes := n.nodes
		for i := 0; i < len(nodes); i++ {
			if nodes[i].id > n.id {
				elect_msg := ELECTMessage(n.id, n.id)
				// send to respective node
				nodes[i].ch_elect <- elect_msg
				fmt.Printf("Node %d sent elect msg to %d.\n", n.id, nodes[i].id)
				logger.Printf("Node %d sent elect msg to %d.\n", n.id, nodes[i].id)
			}
		}
		time.Sleep(n.timeout)
		go n.broadcast_victory()
		return

	}

}

// broadcast victory message
func (n *Node) broadcast_victory() {

	select {
	case msg := <-n.ch_stop_elect:
		fmt.Printf("Node %d %s broadcast_victory.\n", n.id, msg.message)
		logger.Printf("Node %d %s broadcast_victory.\n", n.id, msg.message)
		return
	default:
		// update state
		n.mutex.Lock()
		n.state = BROADCATING
		n.mutex.Unlock()

		nodes := n.nodes

		// send to all node except itself
		for i := 0; i < len(nodes); i++ {
			if nodes[i].id != n.id {
				victory_msg := VICTORYMessage(n.id, n.id)
				nodes[i].ch_elect <- victory_msg
				fmt.Printf("Node %d is broadcasting a victory to %d.\n", n.id, nodes[i].id)
				logger.Printf("Node %d is broadcasting a victory to %d.\n", n.id, nodes[i].id)
			}

		}

		// update role and state
		n.mutex.Lock()
		n.state = NORMAL
		n.role = COORDINATOR
		n.coordinator_id = n.id
		n.mutex.Unlock()

		// switch role
		n.ch_elect <- STOPMessage(n.id)
		n.ch_sync <- STOPMessage(n.id)
		n.ch_role_switch <- STOPMessage(n.id)

		fmt.Printf("Node %d sent role switch messgae.\n", n.id)
		logger.Printf("Node %d sent role switch messgae.\n", n.id)

	}
}

func (n *Node) fail() {
	n.state = DOWN

	fmt.Printf("Node %d failed.\n", n.id)
	logger.Printf("Node %d failed.\n", n.id)
}

func (n *Node) wakeup() {
	// reset own state to RE
	n.state = SELF_ELECTING
	n.role = REPLICA
	go n.self_elect()
}
