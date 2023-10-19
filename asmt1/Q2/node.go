package main

import (
	"fmt"
	"sync"
	"time"
)

type Node struct {
	id             int
	ch_sync        chan Message  //channal listening to SYNC message
	ch_elect       chan Message  //channal listening to ELECT, ACK, VICTORY message
	ch_stop_elect  chan Message  // channal listening to STOP mesage to stop self_elect() or broadcast_victory()
	ch_role_switch chan Message  // channal listening to STOP mesage to switch between Coordinator and worker
	data           Data          // an arbitrary data structure
	sync_interval  time.Duration // coordinator heartbeat time interval
	timeout        time.Duration // 2T(m) + T(p)
	role           string        // coordinator or worker
	nodes          []*Node       // all nodes in the network
	coordinator_id int           // the current coordinator's id in the network
	state          string
	mutex          sync.Mutex
}

func NewNode(id int, sync_interval time.Duration, timeout time.Duration, role string, nodes []*Node) *Node {
	return &Node{
		id:             id,
		ch_sync:        make(chan Message),
		ch_elect:       make(chan Message),
		ch_stop_elect:  make(chan Message),
		ch_role_switch: make(chan Message),
		data:           Data{},
		sync_interval:  sync_interval,
		timeout:        timeout,
		role:           role,
		nodes:          nodes,
		coordinator_id: len(nodes) - 1, //default: highest id
		state:          NORMAL,
		mutex:          sync.Mutex{},
	}
}

// send election message to higher id nodes
func (n *Node) self_elect() {

	fmt.Printf("Node %d started an election.\n", n.id)
	logger.Printf("Node %d started an election.\n", n.id)

	select {

	case msg := <-n.ch_stop_elect:
		// stop self elect
		fmt.Printf("Node %d %s.\n", n.id, msg.message)
		logger.Printf("Node %d %s.\n", n.id, msg.message)
		return

	default:
		// send self elect msg to node with higher ids
		nodes := n.nodes
		for i := 0; i < len(nodes); i++ {
			if nodes[i].id > n.id {
				elect_msg := ELECTMessage(n.id, n.id)
				nodes[i].ch_elect <- elect_msg
				fmt.Printf("Node %d sent elect msg to %d.\n", n.id, nodes[i].id)
				logger.Printf("Node %d sent elect msg to %d.\n", n.id, nodes[i].id)
			}
		}
		// wait for ACK message. If no, broadcast victory
		time.Sleep(n.timeout)
		go n.broadcast_victory()
		return

	}

}

// broadcast victory message
func (n *Node) broadcast_victory() {

	select {
	case msg := <-n.ch_stop_elect:
		// stop broadcasting victory
		fmt.Printf("Node %d %s broadcast_victory.\n", n.id, msg.message)
		logger.Printf("Node %d %s broadcast_victory.\n", n.id, msg.message)
		return

	default:
		// update state
		n.mutex.Lock()
		n.state = BROADCATING
		n.mutex.Unlock()

		nodes := n.nodes

		// send victory to all node except itself
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

// simulate node failure
func (n *Node) fail() {
	n.mutex.Lock()
	n.state = DOWN
	n.mutex.Unlock()

	fmt.Printf("Node %d failed.\n", n.id)
	logger.Printf("Node %d failed.\n", n.id)
}

// simulate node wakeup
func (n *Node) fail_during_election() {
	fmt.Printf("Node %d fail_during_election.\n", n.id)
	logger.Printf("Node %d fail_during_election.\n", n.id)
	n.ch_stop_elect <- STOPMessage(n.id)
	n.mutex.Lock()
	n.state = DOWN
	n.mutex.Unlock()
}

// Worker listen to sync channel
func (n *Node) worker_sync() {

	for {
		select {
		case msg := <-n.ch_sync:
			// if alive, process msg
			if n.state != DOWN {
				fmt.Printf("Node %d received a message %s.\n", n.id, msg.message)
				logger.Printf("Node %d received a message %s.\n", n.id, msg.message)

				switch msg.message_type {
				case STOP:
					return
				case SYNC:
					switch n.state {
					case NORMAL:
						coor_data := msg.message.(Data) // SYNC come with Data in message field
						n.data = coor_data              // update local data
					}
				case SELF_ELECTING, BROADCATING:
					// stops electing if receive sync again
					n.ch_stop_elect <- STOPMessage(n.id)
					n.mutex.Lock()
					n.state = NORMAL
					n.mutex.Unlock()
				}
			}

		case <-time.After(n.timeout):

			//select new coordinator using bully algo
			if n.coordinator_id != n.id && n.state == NORMAL {
				fmt.Printf("Node %d detected coordinator %d failed.\n", n.id, n.coordinator_id)
				logger.Printf("Node %d detected coordinator %d failed.\n", n.id, n.coordinator_id)

				n.mutex.Lock()
				n.state = SELF_ELECTING
				n.mutex.Unlock()

				go n.self_elect()
			}
		}
	}
}

// Worker listen to elect channel
func (n *Node) worker_elect() {

	for {
		msg := <-n.ch_elect

		// if alive, process msg
		if n.state != DOWN {

			fmt.Printf("Node %d received a message %s\n", n.id, msg.message)
			logger.Printf("Node %d received a message %s\n", n.id, msg.message)

			switch msg.message_type {
			case STOP:
				return
			case ELECT:
				switch n.state {
				case NORMAL:
					// haven't realized coor down
					// send ack first, and start election
					if msg.sender_id < n.id {
						// refuse
						ack := ACKMessage(msg.sender_id, n.id)
						n.nodes[msg.sender_id].ch_elect <- ack

						fmt.Printf("Node %d sent ACK to node %d.\n", n.id, msg.sender_id)
						logger.Printf("Node %d sent ACK to node %d.\n", n.id, msg.sender_id)

						if n.coordinator_id != n.id {
							n.mutex.Lock()
							n.state = SELF_ELECTING
							n.mutex.Unlock()
							go n.self_elect() // start election in background
						}
					}
				case SELF_ELECTING:
					if msg.sender_id < n.id {
						// refuse
						ack := ACKMessage(msg.sender_id, n.id)
						n.nodes[msg.sender_id].ch_elect <- ack

						fmt.Printf("Node %d sent ACK to node %d.\n", n.id, msg.sender_id)
						logger.Printf("Node %d sent ACK to node %d.\n", n.id, msg.sender_id)
					}
				}
			case ACK:
				switch n.state {
				case SELF_ELECTING:
					// stop self electing and go back to normal
					n.ch_stop_elect <- STOPMessage(n.id)
					n.mutex.Lock()
					n.state = NORMAL
					n.mutex.Unlock()
				}
			case VICTORY:
				switch n.state {
				case NORMAL:
					if msg.sender_id > n.id {
						// agree on new coordinator
						n.mutex.Lock()
						n.coordinator_id = msg.sender_id
						n.mutex.Unlock()
					} else {
						// self elect
						n.mutex.Lock()
						n.state = SELF_ELECTING
						n.mutex.Unlock()
						go n.self_elect()
					}

				case SELF_ELECTING:
					// stop self electing and ackowledge the new coordinator
					if msg.sender_id > n.id {
						n.ch_stop_elect <- STOPMessage(n.id)
						n.mutex.Lock()
						n.state = NORMAL
						n.coordinator_id = msg.sender_id
						n.mutex.Unlock()
					}
				}
			}
		}

	}
}

// Coordinator listen to elect channel
func (n *Node) coor_elect() {
	for {
		msg := <-n.ch_elect

		// if alive, process msg
		if n.state != DOWN {
			fmt.Printf("Coordinator %d received a message %s.\n", n.id, msg.message)
			logger.Printf("Coordinator %d received a message %s.\n", n.id, msg.message)

			switch msg.message_type {
			case STOP:
				return
			case ELECT:
				switch n.state {
				// when a new coor wakeup
				case NORMAL:
					if msg.sender_id > n.id {
						// step down and switch role
						n.ch_stop_elect <- STOPMessage(n.id)
						n.mutex.Lock()
						n.state = NORMAL
						n.coordinator_id = msg.sender_id
						n.mutex.Unlock()
						return
					} else {
						// send victory again in case it didn't catch
						victory_msg := VICTORYMessage(n.id, n.id)
						n.nodes[msg.sender_id].ch_elect <- victory_msg
					}
				}
			case VICTORY:
				switch n.state {
				case NORMAL:
					if msg.sender_id > n.id {
						// agree on new coordinator
						n.mutex.Lock()
						n.coordinator_id = msg.sender_id
						n.mutex.Unlock()
						// stop current sync
						n.ch_role_switch <- STOPMessage(n.id)
						n.ch_sync <- STOPMessage(n.id)
						n.ch_elect <- STOPMessage(n.id)
					}
				}
			}
		}

	}
}

// Coordinator send sync message to sync channel
func (n *Node) coor_sync() {
	for {
		// if alive, process msg
		if n.state != DOWN {
			msg := SYNCMessage(n.id, n.data)

			// send sync message to all other machines.
			for i := 0; i < len(n.nodes); i++ {
				if n.nodes[i].id < n.id {
					n.nodes[i].ch_sync <- msg
				}
			}
			fmt.Printf("Coordinator %d sent a sync message %s.\n", n.id, msg.message)
			logger.Printf("Coordinator %d sent a sync message %s.\n", n.id, msg.message)

			// sleep for sync_interval
			time.Sleep(n.sync_interval)
		}
	}
}
