package main

import (
	"container/heap"
	"sync"
	"time"
)

type Manager struct {
	id      int // MANAGER_ID(0) ; replica is -1
	ch      chan Message
	sig_ch  chan Message    // for monitor queue
	clock   int             //lamport clock
	servers []*Server       // all servers in the network
	records []ManagerRecord // all page records
	sync.Mutex
}

func NewManager(servers []*Server, records []ManagerRecord) *Manager {
	return &Manager{
		id:      CM_ID,
		ch:      make(chan Message),
		sig_ch:  make(chan Message),
		clock:   0,
		servers: servers,
		records: records,
	}
}

func (m *Manager) listen() {
	for {
		select {
		// receive data
		case msg := <-m.ch:

			m.updateClock(msg.clock)

			if msg.message_type == RD_REQ {

				go m.onReceiveReadReq(msg)
				logger.Printf("Manager received Server %d's read request at clock %d.\n", msg.sender_id, msg.clock)

			} else if msg.message_type == WR_REQ {

				go m.onReceiveWriteReq(msg)
				logger.Printf("Manager received Server %d's write request at clock %d.\n", msg.sender_id, msg.clock)

			} else if msg.message_type == RD_CFM {
				go m.onReceiveReadConfirm(msg)
				logger.Printf("Manager received Server %d's read confirm at clock %d.\n", msg.sender_id, msg.clock)

			} else if msg.message_type == WR_CFM {
				go m.onReceiveWriteConfirm(msg)
				logger.Printf("Manager received Server %d's write confirm at clock %d.\n", msg.sender_id, msg.clock)

			} else if msg.message_type == INV_CFM {
				go m.onReceiveInvConfirm(msg)
				logger.Printf("Manager received Server %d's invalidation confirm at clock %d.\n", msg.sender_id, msg.clock)

			}
		case <-time.After(timeout):
			return
		}
	}

}

func (m *Manager) monitor() {
	for {
		select {
		// receive data
		case msg := <-m.sig_ch:
			m.updateOwnClock()
			logger.Println("Signal channel received msg ", msg, " at clock ", m.clock)

			for i := 0; i < len(m.records); i++ {
				if m.records[i].page_num == msg.page_num {

					m.Lock()
					if !m.records[i].writing {
						m.records[i].writing = true
						// execute req at the head
						go func() {
							msgHead := heap.Pop(&m.records[i].queue).(*Message)
							logger.Printf("Manager poped s%d's wr req for page %d.\n", msgHead.requester, msgHead.page_num)
							if m.isCopyEmpty(msgHead.page_num) {
								m.forwardReq(msgHead.requester, msgHead.page_num, WRITE)
							} else {
								m.requestInvalidateCopy(msgHead.sender_id, msg.page_num)
							}
						}()
					}
					m.Unlock()
					break
				}
			}

		case <-time.After(timeout):
			return
		}
	}

}

// broadcast request to enter cs
func (m *Manager) forwardReq(requester int, requested_page_num int, operation string) {

	m.updateOwnClock()

	for i := 0; i < len(m.records); i++ {
		if m.records[i].page_num == requested_page_num {
			server_id := m.records[i].owner

			logger.Printf("Manager forward %s request for page %d at clock %d.\n", operation, requested_page_num, m.clock)

			msg := ForwardMessage(requester, server_id, requested_page_num, m.clock, operation)
			m.servers[server_id].ch <- msg

			return
		}
	}

	logger.Printf("Page %d is not recorded in manager.\n", requested_page_num)
}

// reply if conditions met
func (m *Manager) requestInvalidateCopy(requester int, requested_page_num int) {

	for i := 0; i < len(m.records); i++ {
		if m.records[i].page_num == requested_page_num {
			for j := 0; j < len(m.records[i].copySet); j++ {
				server_id := m.records[i].copySet[j]
				logger.Printf("Manager request server %d to invalidate cache for page %d at clock %d.\n", server_id, requested_page_num, m.clock)
				msg := InvalidReqMessage(requester, server_id, requested_page_num, m.clock)
				m.servers[server_id].ch <- msg
			}
			return
		}
	}
	logger.Printf("Page %d is not recorded in manager.\n", requested_page_num)
}

func (m *Manager) onReceiveReadReq(msg Message) {
	m.forwardReq(msg.requester, msg.page_num, READ)
}

func (m *Manager) onReceiveReadConfirm(msg Message) {
	m.addToCopySet(msg.page_num, msg.sender_id)
}

func (m *Manager) onReceiveWriteReq(msg Message) {

	for i := 0; i < len(m.records); i++ {
		if m.records[i].page_num == msg.page_num {
			m.Lock()
			heap.Push(&m.records[i].queue, &msg)
			m.Unlock()
			logger.Printf("Manager pushed server %d's wr req for page %d at clock %d.\n", msg.requester, msg.page_num, m.clock)
			m.sig_ch <- NoticeMessage(msg.requester, msg.page_num, msg.clock)
			return
		}
	}
	logger.Printf("Page %d not found in manager.\n", msg.page_num)
}

func (m *Manager) onReceiveInvConfirm(msg Message) {
	m.removeCopy(msg.sender_id)
	if m.isCopyEmpty(msg.page_num) {
		m.forwardReq(msg.requester, msg.page_num, WRITE)
	}
}

func (m *Manager) onReceiveWriteConfirm(msg Message) {
	pn := msg.page_num
	for i := 0; i < len(m.records); i++ {
		if m.records[i].page_num == pn {
			m.records[i].writing = false
			m.records[i].owner = msg.sender_id
			if m.records[i].queue.Len() != 0 {
				m.sig_ch <- NoticeMessage(msg.requester, pn, msg.clock)
			}
		}
	}
}

// update clock
func (m *Manager) updateClock(msgClock int) {
	m.Lock()
	m.clock = max(msgClock, m.clock) + 1
	m.Unlock()

	// logger.Printf("Server %d's clock updated: %d.\n", m.id, m.clock)
}

// update clock
func (m *Manager) updateOwnClock() {
	m.Lock()
	m.clock = m.clock + 1
	m.Unlock()

	// logger.Printf("Server %d's clock updated: %d.\n", m.id, m.clock)
}

func (m *Manager) removeCopy(id int) {
	for i := 0; i < len(m.records); i++ {
		for j := 0; j < len(m.records[i].copySet); j++ {
			if m.records[i].copySet[j] == id {
				m.records[i].copySet = append(m.records[i].copySet[:j], m.records[i].copySet[j+1:]...)
				logger.Printf("Manager remove server %d from page %d copyset.\n", id, m.records[i].page_num)
			}
		}
	}
}

func (m *Manager) isCopyEmpty(page_num int) bool {
	for _, record := range m.records {
		if record.page_num == page_num {
			return len(record.copySet) == 0
		}
	}
	logger.Printf("Page %d not found in manager.\n", page_num)
	return false
}

func (m *Manager) addToCopySet(page_num int, server_id int) {
	m.records[page_num].copySet = append(m.records[page_num].copySet, server_id)
	logger.Printf("Manager adds server %d to page %d's copyset.\n", server_id, page_num)
}

func (m *Manager) start() {
	go m.listen()
	go m.monitor()
}
