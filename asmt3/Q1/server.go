package main

import (
	"sync"
	"time"
)

type Server struct {
	id      int
	ch      chan Message
	clock   int       //lamport clock
	servers []*Server // all servers in the network
	manager *Manager
	records []ServerRecord // record of pages and access
	sync.Mutex
}

func NewServer(id int, servers []*Server, manager *Manager, records []ServerRecord) *Server {
	return &Server{
		id:      id,
		ch:      make(chan Message),
		clock:   0,
		servers: servers,
		manager: manager,
		records: records,
	}
}

func (s *Server) listen() {
	for {
		select {
		// receive data
		case msg := <-s.ch:

			s.updateClock(msg.clock)

			if msg.message_type == RD_FWD {

				go s.onReceiveReadForward(msg)
				logger.Printf("Server %d received Manager's read forward at clock %d.\n", s.id, msg.clock)

			} else if msg.message_type == WR_FWD {

				go s.onReceiveWriteForward(msg)
				logger.Printf("Server %d received Manager's write forward at clock %d.\n", s.id, msg.clock)

			} else if msg.message_type == INV_REQ {

				go s.onReceiveInvReq(msg)
				logger.Printf("Server %d received Manager's invalidation request at clock %d.\n", s.id, msg.clock)

			} else if msg.message_type == SD_RD_PAGE {

				go s.onReceiveReadPage(msg)
				logger.Printf("Server %d received Server %d's sent page at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == SD_WR_PAGE {

				go s.onReceiveWritePage(msg)
				logger.Printf("Server %d received Server %d's sent page at clock %d.\n", s.id, msg.sender_id, msg.clock)

			}
		case <-time.After(timeout):
			return
		}
	}

}

func (s *Server) read(page_num int) {
	logger.Printf("Server %d wants to read page %d...\n", s.id, page_num)

	// if has valid local copy
	for _, record := range s.records {
		if record.page_num == page_num && record.access != NIL {
			logger.Printf("Server %d is reading page %d at clock %d", s.id, page_num, s.clock) // simulate read
			logger.Printf("Server %d finished reading page %d.\n", s.id, page_num)
			return
		}
	}
	// no local copy, request manager
	logger.Printf("Server %d request to read page %d to manager at clock %d.\n", s.id, page_num, s.clock)
	msg := RequestMessage(s.id, page_num, s.clock, READ)
	s.manager.ch <- msg
}

func (s *Server) write(page_num int) {
	logger.Printf("Server %d wants to write page %d...\n", s.id, page_num)

	// check if s is owner
	for _, record := range s.records {
		if record.page_num == page_num && record.access == RW {
			time.Sleep(1) // simulate write
			logger.Printf("Server %d finished writing page %d.\n", s.id, page_num)
			return
		}
	}
	// not owner or does not have a copy
	logger.Printf("Server %d request to write page %d to manager at clock %d.\n", s.id, page_num, s.clock)
	msg := RequestMessage(s.id, page_num, s.clock, WRITE)
	s.manager.ch <- msg

}

func (s *Server) sendPage(page_num int, receiver_id int, operation string) {

	s.updateOwnClock()

	for _, record := range s.records {
		if record.page_num == page_num && record.access == RW {
			msg := SendPageMessage(s.id, receiver_id, record.page, s.clock, operation)
			s.servers[receiver_id].ch <- msg
			logger.Printf("Server %d sent page %d to Server %d at clock %d.\n", s.id, page_num, receiver_id, msg.clock)
			break
		}
	}
}

// broadcast request to enter cs
func (s *Server) confirm(requester int, page_num int, operation string) {

	s.updateOwnClock()

	msg := ConfirmMessage(requester, s.id, page_num, s.clock, operation)
	s.manager.ch <- msg

	logger.Printf("Server %d sent %s confirm for page %d to manager at clock %d.\n", s.id, operation, page_num, msg.clock)
}

func (s *Server) onReceiveReadForward(msg Message) {
	s.sendPage(msg.page_num, msg.requester, READ)
}

func (s *Server) onReceiveWriteForward(msg Message) {
	s.sendPage(msg.page_num, msg.requester, WRITE)
	s.markNilAccess(msg.page_num)
}

func (s *Server) onReceiveReadPage(msg Message) {
	// cache page and mark access as read only
	hasRecord := false
	for i := 0; i < len(s.records); i++ {
		if s.records[i].page_num == msg.page_num {
			s.records[i].access = R
			s.records[i].page = msg.page
			hasRecord = true
		}
	}
	if !hasRecord {
		s.records = append(s.records, NewServerRecord(msg.page_num, R, msg.page))
	}
	logger.Printf("Server %d is reading page %d...", s.id, msg.page_num) // simulate read
	s.confirm(s.id, msg.page_num, READ)

}

func (s *Server) onReceiveWritePage(msg Message) {
	// cache page and mark access as RW
	hasRecord := false
	for i := 0; i < len(s.records); i++ {
		if s.records[i].page_num == msg.page_num {
			s.records[i].access = RW
			s.records[i].page = msg.page
			hasRecord = true
		}
	}
	if !hasRecord {
		s.records = append(s.records, NewServerRecord(msg.page_num, RW, msg.page))
	}
	logger.Printf("Server %d is writing page %d at clock %d.", s.id, msg.page_num, s.clock) // simulate write
	s.confirm(s.id, msg.page_num, WRITE)
}

func (s *Server) onReceiveInvReq(msg Message) {
	s.markNilAccess(msg.page_num)
}

// update clock
func (s *Server) updateClock(msgClock int) {
	s.Lock()
	s.clock = max(msgClock, s.clock) + 1
	s.Unlock()

	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}

// update clock
func (s *Server) updateOwnClock() {
	s.Lock()
	s.clock = s.clock + 1
	s.Unlock()

	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}

func (s *Server) markNilAccess(page_num int) {
	// remove access
	for _, record := range s.records {
		if record.page_num == page_num {
			record.access = NIL
		}
	}
}
