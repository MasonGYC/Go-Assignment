package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

const (
	// server states
	WAITING_FOR_REPLY string = "WAITING_FOR_REPLY"
	NORMAL            string = "NORMAL"
)

type Server struct {
	id            int
	ch            chan Message
	queue         RequestPriorityQueue
	clock         int       //lamport clock
	servers       []*Server // all servers in the network
	state         string    // REQUESTING, NORMAL
	reply_counter int       // number of reply to be received
	holding_reply []Request // number of holding replies
	sync.Mutex
}

func NewServer(id int, num_servers int, servers []*Server) *Server {
	pq := make(RequestPriorityQueue, 0)
	heap.Init(&pq)
	return &Server{
		id:            id,
		ch:            make(chan Message),
		queue:         pq,
		clock:         0,
		servers:       servers,
		state:         NORMAL,
		reply_counter: len(servers) - 1,
		holding_reply: make([]Request, 0),
	}
}

func (s *Server) listen() {
	for {
		// receive data
		msg := <-s.ch

		s.updateClock(msg.clock)

		if msg.message_type == REQ {

			go s.onReceiveReq(msg)
			fmt.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)
			logger.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)

		} else if msg.message_type == REP {

			go s.onReceiveRep(msg)
			fmt.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)
			logger.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)

		}
	}
}

// broadcast request to enter cs
func (s *Server) request() {
	if s.state == NORMAL {

		s.updateOwnClock()

		s.Lock()
		s.state = WAITING_FOR_REPLY
		s.reply_counter = len(s.servers) - 1
		s.Unlock()

		// add to queue
		req := &Request{
			value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", s.id, s.id),
			clock:     s.clock,
			requester: s.id,
		}
		heap.Push(&s.queue, req)

		// broadcast requests
		servers := s.servers
		for i := 0; i < len(servers); i++ {
			if servers[i].id != s.id {
				req_msg := REQMessage(s.id, servers[i].id, req.clock, s.clock)
				servers[i].ch <- req_msg

				fmt.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, s.clock)
				logger.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, s.clock)
			}
		}
	} else {
		fmt.Printf("Server %d has already requested.\n", s.id)
		logger.Printf("Server %d has already requested.\n", s.id)

	}
}

// reply if conditions met
func (s *Server) reply(requester_id int, request_clock int) {

	s.updateOwnClock()

	// send reply
	res_msg := REPMessage(s.id, requester_id, request_clock, s.clock)
	s.servers[requester_id].ch <- res_msg

	fmt.Printf("Server %d replys to %d at %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d replys to %d at %d.\n", s.id, requester_id, s.clock)
}

// empty_queue locks to cs
func (s *Server) empty_queue() {

	s.updateOwnClock()

	// empty queue by replying to all other requests
	for {
		request := s.queue.Peek()

		if request != nil {
			request := heap.Pop(&s.queue).(*Request)
			if request.requester != s.id {
				s.reply(request.requester, request.clock)
			}

			fmt.Printf("Server %d poped request from %d and replied at clock %d.\n", s.id, request.requester, s.clock)
			logger.Printf("Server %d Poped reuqest from %d.\n", s.id, request.requester)

		} else {
			break
		}
	}
}

// simulate execution of critical section
func (s *Server) executeCriticalSection() {
	s.updateOwnClock()
	time.Sleep(2 * time.Second)
	fmt.Printf("Server %d has finished cs execution.\n", s.id)
	logger.Printf("Server %d has finished cs execution.\n", s.id)
}

func (s *Server) onReceiveReq(msg Message) {

	// Check whether any reply is pending for an earlier request req(T’) in Qi

	req_at_head := s.queue.Peek()
	if req_at_head != nil {
		fmt.Printf("Server %d has req from server %d at clock %d at head of queue.\n", s.id, req_at_head.requester, req_at_head.clock)
		logger.Printf("Server %d has req from server %d at clock %d at head of queue.\n", s.id, req_at_head.requester, req_at_head.clock)

		if req_at_head.clock < msg.clock || (req_at_head.clock == msg.clock && req_at_head.requester > msg.sender_id) {
			// add to queue
			req := &Request{
				value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", msg.sender_id, s.id),
				clock:     msg.clock,
				requester: msg.sender_id,
			}
			heap.Push(&s.queue, req)
			fmt.Printf("Server %d has pushed req from %d to queue .\n", s.id, req.requester)
			logger.Printf("Server %d has pushed req from %d to queue .\n", s.id, req.requester)

		} else {
			s.reply(msg.sender_id, msg.clock)
		}

	} else {
		s.reply(msg.sender_id, msg.clock)
	}

}
func (s *Server) onReceiveRep(msg Message) {

	if s.state == WAITING_FOR_REPLY {

		s.Lock()
		s.reply_counter--
		s.Unlock()

		// check replies

		if s.reply_counter == 0 {
			// change state to normal
			s.Lock()
			s.state = NORMAL
			s.reply_counter = len(s.servers) - 1
			s.Unlock()

			s.executeCriticalSection()
			s.empty_queue()

		}
	}
}

// update clock
func (s *Server) updateClock(msgClock int) {
	s.Lock()
	s.clock = max(msgClock, s.clock) + 1
	s.Unlock()

	// fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}

// update clock
func (s *Server) updateOwnClock() {
	s.Lock()
	s.clock = s.clock + 1
	s.Unlock()

	// fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}
