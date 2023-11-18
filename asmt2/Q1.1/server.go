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
	fmt.Printf("Server %d has pushed req from %d to queue.\n", s.id, req.requester)
	logger.Printf("Server %d has pushed req from %d to queue.\n", s.id, req.requester)

	// If waiting for REPLY from j for an earlier request T, wait until j replies to you
	req_at_head := s.queue.Peek()
	fmt.Printf("Server %d has req from server %d at clock %d at head of queue.\n", s.id, req_at_head.requester, req_at_head.clock)
	logger.Printf("Server %d has req from server %d at clock %d at head of queue.\n", s.id, req_at_head.requester, req_at_head.clock)

	if req_at_head.requester == s.id {
		// hold the reply
		s.holding_reply = append(s.holding_reply, *req)
		fmt.Printf("Server %d holds reply to %d.\n", s.id, req.requester)
		logger.Printf("Server %d holds reply to %d.\n", s.id, req.requester)
	} else {
		s.reply(msg.sender_id, msg.clock)
	}

}
func (s *Server) onReceiveRep(msg Message) {

	if s.state == WAITING_FOR_REPLY {

		s.Lock()
		s.reply_counter--
		s.Unlock()

		// check replies and whether at the head of queue
		req := s.queue.Peek()
		fmt.Printf("Server %d's reply_counter is %d.\n", s.id, s.reply_counter)
		logger.Printf("Server %d's reply_counter is %d.\n", s.id, s.reply_counter)

		if req != nil {
			if s.reply_counter == 0 {

				// change state to normal
				s.Lock()
				s.state = NORMAL
				s.reply_counter = len(s.servers) - 1
				s.Unlock()

				go func() {
					// clear holding reply
					for i := 0; i < len(s.holding_reply); i++ {
						req := s.holding_reply[i]
						fmt.Printf("Server %d clear holding reply from %d at %d.\n", s.id, req.requester, req.clock)
						logger.Printf("Server %d clear holding reply from %d at %d.\n", s.id, req.requester, req.clock)

						s.reply(req.requester, req.clock)
					}
					s.holding_reply = make([]Request, 0)
				}()

				s.executeCriticalSection()
				s.release()

			}
		}
	}
}

func (s *Server) onReceiveRls() {
	// pop queue
	heap.Pop(&s.queue)
	fmt.Printf("Server %d has poped req from queue.\n", s.id)
	logger.Printf("Server %d has poped req from queue.\n", s.id)

	// check replies and whether at the head of queue
	if s.state == WAITING_FOR_REPLY && s.reply_counter == 0 {
		req := s.queue.Peek()
		if req != nil && req.requester == s.id {
			// change state to normal
			s.Lock()
			s.state = NORMAL
			s.reply_counter = len(s.servers) - 1
			s.Unlock()

			s.executeCriticalSection()
			s.release()
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
