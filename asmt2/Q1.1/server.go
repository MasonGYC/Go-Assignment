package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	id                         int
	ch                         chan Message
	queue                      RequestPriorityQueue
	clock                      int       //lamport clock
	servers                    []*Server // all servers in the network
	reply_counter              int       // number of reply to be received
	holding_reply              []Request // number of holding replies
	waiting_for_reply_at_clock int       // the clock of the request that waits for reply
	sync.Mutex
}

func NewServer(id int, num_servers int, servers []*Server) *Server {
	pq := make(RequestPriorityQueue, 0)
	heap.Init(&pq)
	return &Server{
		id:                         id,
		ch:                         make(chan Message),
		queue:                      pq,
		clock:                      0,
		servers:                    servers,
		reply_counter:              len(servers) - 1,
		holding_reply:              make([]Request, 0),
		waiting_for_reply_at_clock: -1,
	}
}

func (s *Server) listen() {
	for {
		select {
		// receive data
		case msg := <-s.ch:

			s.updateClock(msg.clock)

			if msg.message_type == REQ {

				go s.onReceiveReq(msg)
				fmt.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == REP {

				go s.onReceiveRep(msg)
				fmt.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == RLS {

				go s.onReceiveRls()
				fmt.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)

			}
		case <-time.After(timeout):
			return
		}
	}

}

// broadcast request to enter cs
func (s *Server) request() {
	if s.waiting_for_reply_at_clock == -1 {

		s.updateOwnClock()

		req := &Request{
			value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", s.id, s.id),
			clock:     s.clock,
			requester: s.id,
		}

		s.Lock()
		s.waiting_for_reply_at_clock = req.clock
		s.reply_counter = len(s.servers) - 1
		s.Unlock()

		// if only 1 server
		if s.reply_counter == 0 {
			s.executeCriticalSection()
		}

		// add to queue
		heap.Push(&s.queue, req)

		fmt.Printf("Server %d make request at clock %d.\n", s.id, req.clock)
		logger.Printf("Server %d make request at clock %d.\n", s.id, req.clock)

		// broadcast requests
		servers := s.servers
		for i := 0; i < len(servers); i++ {
			if servers[i].id != s.id {
				req_msg := REQMessage(s.id, servers[i].id, req.clock, s.clock)
				servers[i].ch <- req_msg

				fmt.Printf("Server %d requests to %d at clock %d.\n", s.id, servers[i].id, s.clock)
				logger.Printf("Server %d requests to %d at clock %d.\n", s.id, servers[i].id, s.clock)
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

	fmt.Printf("Server %d replys to %d at clock %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d replys to %d at clock %d.\n", s.id, requester_id, s.clock)
}

// release locks to cs
func (s *Server) release() {

	s.updateOwnClock()

	// Pop head of Q
	request := heap.Pop(&s.queue).(*Request)
	fmt.Printf("Poped reuqest %d:%s at clock %d.\n", request.clock, request.value, s.clock)

	// broadcast release
	servers := s.servers
	for i := 0; i < len(servers); i++ {
		if servers[i].id != s.id {
			rls_msg := RLSMessage(s.id, servers[i].id, request.clock, s.clock)
			servers[i].ch <- rls_msg

			fmt.Printf("Server %d release to %d at clock %d.\n", s.id, servers[i].id, s.clock)
			logger.Printf("Server %d release to %d at clock %d.\n", s.id, servers[i].id, s.clock)
		}
	}
}

// simulate execution of critical section
func (s *Server) executeCriticalSection() {
	s.updateOwnClock()
	time.Sleep(2 * time.Second)

	time_mutex.Lock()
	now := time.Now()
	// if end_time before now
	if end_time.Compare(now) == -1 {
		end_time = now
		fmt.Print(end_time)
		logger.Print(end_time)
	}
	time_mutex.Unlock()

	fmt.Printf("Server %d has finished cs execution.\n", s.id)
	logger.Printf("Server %d has finished cs execution.\n", s.id)
}

func (s *Server) onReceiveReq(msg Message) {
	// add to queue
	req := &Request{
		value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", msg.sender_id, s.id),
		clock:     msg.message.clock,
		requester: msg.message.requester,
	}
	heap.Push(&s.queue, req)
	fmt.Printf("Server %d pushed req from %d to queue.\n", s.id, req.requester)
	logger.Printf("Server %d pushed req from %d to queue.\n", s.id, req.requester)

	// If waiting for REPLY from j for an earlier request T, wait until j replies to you
	if s.waiting_for_reply_at_clock != -1 {
		if s.waiting_for_reply_at_clock < req.clock || (s.waiting_for_reply_at_clock == req.clock && s.id > req.requester) {
			// hold the reply
			s.holding_reply = append(s.holding_reply, *req)
			fmt.Printf("Server %d holds reply to %d.\n", s.id, req.requester)
			logger.Printf("Server %d holds reply to %d.\n", s.id, req.requester)
		} else {
			s.reply(msg.sender_id, msg.clock)
		}
	} else {
		s.reply(msg.sender_id, msg.clock)
	}

}
func (s *Server) onReceiveRep(msg Message) {

	if s.waiting_for_reply_at_clock != -1 {

		s.Lock()
		s.reply_counter--
		s.Unlock()

		// check replies and whether at the head of queue
		req := s.queue.Peek()
		fmt.Printf("Server %d's reply_counter is %d.\n", s.id, s.reply_counter)
		logger.Printf("Server %d's reply_counter is %d.\n", s.id, s.reply_counter)

		if req != nil {
			if s.reply_counter == 0 && req.requester == s.id {
				s.enterCS()
			}
		}
	}
}

func (s *Server) onReceiveRls() {
	// pop queue
	heap.Pop(&s.queue)
	fmt.Printf("Server %d has poped req from queue.\n", s.id)
	logger.Printf("Server %d has poped req from queue.\n", s.id)

	logger.Printf("Server %d waiting_for_reply_at_clock is: %d\n", s.id, s.waiting_for_reply_at_clock)
	logger.Printf("Server %d reply_counter is: %d\n", s.id, s.reply_counter)

	// check replies and whether at the head of queue
	if s.waiting_for_reply_at_clock != -1 && s.reply_counter == 0 {
		req := s.queue.Peek()
		if req != nil && req.requester == s.id {
			s.enterCS()
		}
	}
}

func (s *Server) enterCS() {
	// change state to normal
	s.Lock()
	s.waiting_for_reply_at_clock = -1
	s.reply_counter = len(s.servers) - 1
	s.Unlock()

	// clear holding reply
	for i := 0; i < len(s.holding_reply); i++ {
		req := s.holding_reply[i]
		fmt.Printf("Server %d clear holding reply from %d at %d.\n", s.id, req.requester, req.clock)
		logger.Printf("Server %d clear holding reply from %d at %d.\n", s.id, req.requester, req.clock)

		s.reply(req.requester, req.clock)
	}
	s.holding_reply = make([]Request, 0)

	s.executeCriticalSection()
	s.release()
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
