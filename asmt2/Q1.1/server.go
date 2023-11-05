package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

const (
	// server states
	REQUESTING string = "REQUESTING"
	NORMAL     string = "NORMAL"
)

type Server struct {
	id            int
	ch            chan Message
	queue         RequestPriorityQueue
	clock         []int
	servers       []*Server // all servers in the network
	reply_counter int       // number of reply to be recived
	state         string    // REQUESTING, NORMAL
	sync.Mutex
}

func NewServer(id int, num_servers int, servers []*Server) *Server {
	pq := make(RequestPriorityQueue, 0)
	heap.Init(&pq)
	return &Server{
		id:            id,
		ch:            make(chan Message),
		queue:         pq,
		clock:         make([]int, num_servers),
		servers:       servers,
		reply_counter: 0,
		state:         NORMAL,
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

		} else if msg.message_type == RES {

			go s.onReceiveRes()
			fmt.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)
			logger.Printf("Server %d received Server %d's reply at clock %d.\n", s.id, msg.sender_id, msg.clock)

		} else if msg.message_type == RLS {

			go s.onReceiveRls()
			fmt.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)
			logger.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)

		}
	}
}

// broadcast request to enter cs
func (s *Server) request() {
	if s.state == NORMAL {

		s.updateOwnClock()

		s.state = REQUESTING
		s.reply_counter = len(s.servers) - 1

		// add to queue
		req := &Request{
			value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", s.id, s.id),
			clock:     s.clock,
			sender_id: s.id,
		}
		heap.Push(&s.queue, req)

		// broadcast requests
		servers := s.servers
		for i := 0; i < len(servers); i++ {
			if servers[i].id != s.id {
				req_msg := REQMessage(s.id, servers[i].id, s.clock)
				servers[i].ch <- req_msg

				fmt.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, s.clock)
				logger.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, s.clock)
			}
		}
	}
}

// reply if conditions met
func (s *Server) reply(requestMsg Message) {

	s.updateOwnClock()

	res_msg := RESMessage(s.id, requestMsg.sender_id, s.clock)
	s.servers[requestMsg.sender_id].ch <- res_msg

	fmt.Printf("Server %d replys to %d at %d.\n", s.id, requestMsg.sender_id, s.clock)
	logger.Printf("Server %d replys to %d at %d.\n", s.id, requestMsg.sender_id, s.clock)
}

// release locks to cs
func (s *Server) release() {

	s.updateOwnClock()

	// Pop head of Q
	item := heap.Pop(&s.queue).(*Request)
	fmt.Printf("Poped reuqest %.2d:%s ", item.clock, item.value)

	// broadcast release
	servers := s.servers
	for i := 0; i < len(servers); i++ {
		if servers[i].id != s.id {
			rls_msg := RLSMessage(s.id, servers[i].id, s.clock)
			servers[i].ch <- rls_msg

			fmt.Printf("Server %d release to %d at %d.\n", s.id, servers[i].id, s.clock)
			logger.Printf("Server %d release to %d at %d.\n", s.id, servers[i].id, s.clock)
		}
	}

	// change state to normal
	s.state = NORMAL
}

// simulate execution of critical section
func (s *Server) executeCriticalSection() {
	s.updateOwnClock()
	time.Sleep(1 * time.Second)
	fmt.Printf("Server %d has finished cs execution.\n", s.id)
	logger.Printf("Server %d has finished cs execution.\n", s.id)
	go s.release()
}

func (s *Server) onReceiveReq(msg Message) {
	// add to queue
	req := &Request{
		value:     msg.message,
		clock:     msg.clock,
		sender_id: msg.sender_id,
	}
	heap.Push(&s.queue, req)

	// If waiting for REPLY from j for an earlier request T, wait until j replies to you
	existing_req := s.queue.Peek()
	if existing_req != nil && s.reply_counter > 0 {
		// hold the reply
	} else {
		s.reply(msg)
	}
}
func (s *Server) onReceiveRes() {
	// check total reply
	if s.state == REQUESTING {

		s.Lock()
		s.reply_counter = s.reply_counter - 1
		s.Unlock()

		if s.reply_counter == 0 {
			// received all replies
			go s.executeCriticalSection()
		}
	}
}

func (s *Server) onReceiveRls() {
	// pop queue
	request := heap.Pop(&s.queue).(*Request)
	fmt.Printf("Server %d poped request from Server %d.\n", s.id, request.sender_id)
	logger.Printf("Server %d poped request from Server %d.\n", s.id, request.sender_id)
}

// update clock
func (s *Server) updateClock(msgClock []int) {
	s.Lock()
	for i := 0; i < len(s.clock); i++ {
		s.clock[i] = max(msgClock[i], s.clock[i])
	}
	s.clock[s.id] = s.clock[s.id] + 1
	s.Unlock()

	fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}

// update clock
func (s *Server) updateOwnClock() {
	s.Lock()
	s.clock[s.id] = s.clock[s.id] + 1
	s.Unlock()

	fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}
