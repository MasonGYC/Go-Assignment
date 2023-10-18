package main

import (
	"fmt"
	"time"
)

func (n *Node) rep_sync() {
	// fmt.Printf("Node %d rep_listen_ch_sync.\n", n.id)
	// logger.Printf("Node %d rep_listen_ch_sync.\n", n.id)

	for {
		select {
		case msg := <-n.ch_sync:
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
		case <-time.After(n.timeout):

			fmt.Printf("Node %d sync timeout\n", n.id)
			logger.Printf("Node %d sync timeout\n", n.id)

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

func (n *Node) rep_elect() {

	// fmt.Printf("Node %d rep_listen_ch_elect.\n", n.id)
	// logger.Printf("Node %d rep_listen_ch_elect.\n", n.id)

	for {
		msg := <-n.ch_elect
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

					// TODO: check correctness
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
				} else {
					// maintain SELF-ELECTING state, wait for VIC/self election succeed
					// do nothing
				}
			}
		case ACK:
			switch n.state {
			case SELF_ELECTING:
				n.ch_stop_elect <- STOPMessage(n.id)
				n.mutex.Lock()
				n.state = NORMAL
				n.mutex.Unlock()
			}
		case VICTORY:
			switch n.state {
			case NORMAL:
				if msg.sender_id > n.id {
					n.mutex.Lock()
					n.coordinator_id = msg.sender_id // agree on new coordinator
					n.mutex.Unlock()
				} else {
					// may not happen but put here in case
					n.mutex.Lock()
					n.state = SELF_ELECTING
					n.mutex.Unlock()
					go n.self_elect() // start election in background
				}

			case SELF_ELECTING:
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

func (n *Node) coor_elect() {
	for {
		msg := <-n.ch_elect
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
					n.ch_stop_elect <- STOPMessage(n.id)
					n.mutex.Lock()
					n.state = NORMAL
					n.coordinator_id = msg.sender_id
					n.mutex.Unlock()
					return // switch role
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

func (n *Node) coor_sync() {
	for {
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
