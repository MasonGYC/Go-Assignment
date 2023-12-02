package main

import (
	"container/heap"
	"sync"
	"time"
)

type Manager struct {
	id             int
	ch             chan Message
	sig_ch         chan Message    // for monitor queue
	hbt_ch         chan Message    //heartbeat
	data_ch        chan Message    // for data update
	clock          int             //lamport clock
	servers        []*Server       // all servers in the network
	records        []ManagerRecord // all page records
	role           string          // unmodifierable
	isInCharge     bool            // is performing primary tasks
	isAlive        bool
	backupManager  *Manager
	primaryManager *Manager
	sync.Mutex
}

func NewManager(id int, servers []*Server, records []ManagerRecord, role string) *Manager {
	var isInCharge bool
	if role == PRIMARY {
		isInCharge = true
	} else {
		isInCharge = false
	}
	return &Manager{
		id:         id,
		ch:         make(chan Message),
		sig_ch:     make(chan Message),
		hbt_ch:     make(chan Message),
		data_ch:    make(chan Message),
		clock:      0,
		servers:    servers,
		records:    records,
		role:       role,
		isAlive:    false,
		isInCharge: isInCharge,
	}
}

func (m *Manager) listenToServer() {
	for {
		select {
		// receive data
		case msg := <-m.ch:
			logger.Printf("Manager %d received Server %d's message at clock %d.\n", m.id, msg.sender_id, msg.clock)

			if m.isAlive {
				m.updateClock(msg.clock)

				if msg.message_type == RD_REQ {

					go m.onReceiveReadReq(msg)
					logger.Printf("Manager %d received Server %d's read request at clock %d.\n", m.id, msg.sender_id, msg.clock)

				} else if msg.message_type == WR_REQ {

					go m.onReceiveWriteReq(msg)
					logger.Printf("Manager %d received Server %d's write request at clock %d.\n", m.id, msg.sender_id, msg.clock)

				} else if msg.message_type == RD_CFM {
					go m.onReceiveReadConfirm(msg)
					logger.Printf("Manager %d received Server %d's read confirm at clock %d.\n", m.id, msg.sender_id, msg.clock)

				} else if msg.message_type == WR_CFM {
					go m.onReceiveWriteConfirm(msg)
					logger.Printf("Manager %d received Server %d's write confirm at clock %d.\n", m.id, msg.sender_id, msg.clock)

				} else if msg.message_type == INV_CFM {
					go m.onReceiveInvConfirm(msg)
					logger.Printf("Manager %d received Server %d's invalidation confirm at clock %d.\n", m.id, msg.sender_id, msg.clock)

				}
			}

		case <-time.After(timeout):
			return
		}
	}

}

func (m *Manager) listenToHeartbeat() {
	for {
		select {
		// receive data
		case msg := <-m.hbt_ch:
			if m.isAlive {
				go m.onReceiveHeartbeat(msg)
			}
		case <-time.After(heartbeat_timeout):
			// the other manager is down
			logger.Printf("Manager %d detects the other manager is down.\n", m.id)
			go func() {
				if !m.isInCharge {
					m.Lock()
					m.isInCharge = true
					m.Unlock()
					for i := 0; i < len(m.servers); i++ {
						logger.Printf("Declare manager %d to be pri to server %d at clock %d.\n", m.id, m.servers[i].id, m.clock)
						go func(i int) {
							m.servers[i].ch <- DeclarePrimaryMessage(m.id, m.servers[i].id, m.clock)
						}(i)
					}
				}
			}()
		}
	}
}

func (m *Manager) listenToData() {
	for {
		msg := <-m.data_ch
		if m.isAlive {
			m.records = msg.records
			logger.Printf("Manager %d updated records.\n", m.id)
		}
	}
}

// primary send heartbeat to backup
func (m *Manager) heartbeat() {
	for {
		if m.isAlive {
			m.updateOwnClock()
			m.sendHeartbeat(m.backupManager)
			logger.Printf("Manager %d sent heartbeat to %d at clock %d.\n", m.id, m.backupManager.id, m.clock)
			time.Sleep(heartbeat_interval)
		}
	}
}

func (m *Manager) sendHeartbeat(toManager *Manager) {
	m.updateOwnClock()
	toManager.hbt_ch <- HeartbeatMessage(m.id, toManager.id, m.clock, m)
}

func (m *Manager) updateData() {
	m.updateOwnClock()
	var toManager *Manager
	if m.role == PRIMARY {
		toManager = m.backupManager
	} else {
		toManager = m.primaryManager
	}
	toManager.data_ch <- RecordUpdateMessage(m.id, toManager.id, m.records, m.clock)

}

func (m *Manager) listenToSignal() {
	for {
		select {
		// receive data
		case msg := <-m.sig_ch:
			m.updateOwnClock()
			logger.Printf("Signal channel received msg from requester %d at clock %d", msg.requester, m.clock)

			for i := 0; i < len(m.records); i++ {
				if m.records[i].page_num == msg.page_num {

					m.Lock()
					logger.Println("record.writing is ", m.records[i].writing, " at clock ", m.clock)
					if !m.records[i].writing {
						m.records[i].writing = true
						// execute req at the head
						go func() {
							msgHead := heap.Pop(&m.records[i].queue).(*Message)
							logger.Printf("Manager %d poped s%d's wr req for page %d.\n", m.id, msgHead.requester, msgHead.page_num)
							go m.updateData()
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

			logger.Printf("Manager %d forward %s request for page %d at clock %d.\n", m.id, operation, requested_page_num, m.clock)

			msg := ForwardMessage(requester, m.id, server_id, requested_page_num, m.clock, operation)
			m.servers[server_id].ch <- msg

			return
		}
	}

	logger.Printf("Page %d is not recorded in manager %d.\n", requested_page_num, m.id)
}

// reply if conditions met
func (m *Manager) requestInvalidateCopy(requester int, requested_page_num int) {

	m.updateOwnClock()

	for i := 0; i < len(m.records); i++ {
		if m.records[i].page_num == requested_page_num {
			for j := 0; j < len(m.records[i].copySet); j++ {
				server_id := m.records[i].copySet[j]
				logger.Printf("Manager %d request server %d to invalidate cache for page %d at clock %d.\n", m.id, server_id, requested_page_num, m.clock)
				msg := InvalidReqMessage(requester, m.id, server_id, requested_page_num, m.clock)
				m.servers[server_id].ch <- msg
			}
			return
		}
	}
	logger.Printf("Page %d is not recorded in manager.\n", requested_page_num)
}

func (m *Manager) onReceiveHeartbeat(msg Message) {
	logger.Printf("Manager %d received heartbeat from %d at clock %d.\n", m.id, msg.sender_id, m.clock)

	if m.isBackup() {
		if !m.isInCharge {
			m.sendHeartbeat(m.primaryManager)
		} else {
			// primary rejoin
			m.Lock()
			m.isInCharge = false
			m.Unlock()
			go m.updateData()
		}
	}
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
			go m.updateData()
			logger.Printf("Manager %d pushed server %d's wr req for page %d at clock %d.\n", m.id, msg.requester, msg.page_num, m.clock)
			m.sig_ch <- NoticeMessage(msg.requester, m.id, msg.page_num, msg.clock)
			return
		}
	}
	logger.Printf("Page %d not found in manager %d.\n", msg.page_num, m.id)
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
				m.sig_ch <- NoticeMessage(msg.requester, m.id, pn, msg.clock)
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
				logger.Printf("Manager %d remove server %d from page %d copyset.\n", m.id, id, m.records[i].page_num)
			}
		}
	}
	go m.updateData()
}

func (m *Manager) isCopyEmpty(page_num int) bool {
	for _, record := range m.records {
		if record.page_num == page_num {
			return len(record.copySet) == 0
		}
	}
	logger.Printf("Page %d not found in manager %d.\n", page_num, m.id)
	return false
}

func (m *Manager) addToCopySet(page_num int, server_id int) {
	m.records[page_num].copySet = append(m.records[page_num].copySet, server_id)
	go m.updateData()
	logger.Printf("Manager %d adds server %d to page %d's copyset.\n", m.id, server_id, page_num)
}

func (m *Manager) start() {
	m.Lock()
	m.isAlive = true
	m.Unlock()
	go m.listenToServer()
	go m.listenToSignal()
	go m.listenToHeartbeat()
	go m.listenToData()
	if !m.isBackup() {
		go m.heartbeat()
	}
	logger.Printf("Manager %d started.\n", m.id)
}

func (m *Manager) down() {

	m.Lock()
	m.isAlive = false
	m.clock = 0
	m.records = make([]ManagerRecord, 0)
	m.Unlock()

	logger.Printf("Manager %d is down.\n", m.id)
}

func (m *Manager) rejoin() {
	logger.Printf("Manager %d rejoined.\n", m.id)
	m.Lock()
	m.isAlive = true
	m.isInCharge = true
	m.Unlock()
	for i := 0; i < len(m.servers); i++ {
		logger.Printf("Declare manager %d to be pri to server %d at clock %d.\n", m.id, m.servers[i].id, m.clock)
		go func(i int) {
			m.servers[i].ch <- DeclarePrimaryMessage(m.id, m.servers[i].id, m.clock)
		}(i)
	}
}

func (m *Manager) isBackup() bool {
	return m.role == BACKUP
}
