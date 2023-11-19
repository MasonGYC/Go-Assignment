package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Server struct {
	id                  int
	ch                  chan Message
	clock               int
	servers             []*Server            // all servers in the network
	voted_req           *Request             // voted request
	requesting          bool                 // true if requested
	received_votes      []int                // id of voters; len() = num of votes received
	received_rescinds   []int                // record rescind voters received in case vote arrives later than rescind
	pending_votes_queue RequestPriorityQueue // votes to be sent
	executing_cs        bool                 // is executing cs
	sync.Mutex
}

func NewServer(id int, num_servers int, servers []*Server) *Server {
	pq := make(RequestPriorityQueue, 0)
	heap.Init(&pq)
	return &Server{
		id:                  id,
		ch:                  make(chan Message),
		clock:               0,
		servers:             servers,
		requesting:          false,
		received_votes:      make([]int, 0),
		received_rescinds:   make([]int, 0),
		pending_votes_queue: pq,
		executing_cs:        false,
		voted_req:           nil,
	}
}

func (s *Server) hasMajorityVotes() bool {
	received_votes := len(s.received_votes)
	num_of_servers := len(s.servers)
	if num_of_servers%2 == 1 {
		return received_votes > (num_of_servers-1)/2
	} else {
		return received_votes > num_of_servers/2
	}

}

// listen to messages
func (s *Server) listen() {
	for {
		select {
		case msg := <-s.ch:
			s.updateClock(msg.clock)

			if msg.message_type == REQ {

				go s.onReceiveReq(msg)
				fmt.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's request at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == VOTE {

				go s.onReceiveVote(msg.sender_id)
				fmt.Printf("Server %d received Server %d's vote at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's vote at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == RLS {

				go s.onReceiveRls(msg.sender_id, RLS)
				fmt.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == RESCIND_RLS {

				go s.onReceiveRls(msg.sender_id, RESCIND_RLS)
				fmt.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's release at clock %d.\n", s.id, msg.sender_id, msg.clock)

			} else if msg.message_type == RESCIND {

				go s.onReceiveRcd(msg.sender_id)
				fmt.Printf("Server %d received Server %d's rescind at clock %d.\n", s.id, msg.sender_id, msg.clock)
				logger.Printf("Server %d received Server %d's rescind at clock %d.\n", s.id, msg.sender_id, msg.clock)

			}
		case <-time.After(timeout):
			return

		}

	}
}

// broadcast request to enter cs
func (s *Server) request() {

	s.updateOwnClock()

	servers := s.servers
	req := &Request{
		value:     fmt.Sprintf("Request from server %d to server %d to enter critical section.", s.id, s.id),
		clock:     s.clock,
		requester: s.id,
	}
	fmt.Printf("Server %d made request at %d.\n", s.id, req.clock)
	logger.Printf("Server %d made request at %d.\n", s.id, req.clock)

	s.Lock()
	s.requesting = true
	s.Unlock()

	// if only one server
	if len(servers) == 1 {
		s.executeCriticalSection()
	}

	for i := 0; i < len(servers); i++ {
		req_msg := REQMessage(s.id, servers[i].id, req.clock, s.clock)

		fmt.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, req.clock)
		logger.Printf("Server %d requests to %d at %d.\n", s.id, servers[i].id, req.clock)
		servers[i].ch <- req_msg
	}

}

// vote
func (s *Server) vote(requester_id int, request_clock int) {

	s.updateOwnClock()

	req := &Request{
		clock:     request_clock,
		requester: requester_id,
	}

	s.Lock()
	s.voted_req = req
	s.Unlock()

	vote_msg := VOTEMessage(s.id, requester_id, request_clock, s.clock)
	fmt.Printf("Server %d vote to %d at %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d vote to %d at %d.\n", s.id, requester_id, s.clock)
	s.servers[requester_id].ch <- vote_msg
}

// rescind vote
func (s *Server) rescindVote(requester_id int) {

	s.updateOwnClock()

	rcd_msg := RESCINDMessage(s.id, requester_id, s.voted_req.clock, s.clock)
	fmt.Printf("Server %d rescind vote to %d at %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d rescind vote to %d at %d.\n", s.id, requester_id, s.clock)
	go func() {
		s.servers[requester_id].ch <- rcd_msg
	}()
}

// release vote from voter
func (s *Server) releaseVote(voter_id int, message_type string) {

	s.updateOwnClock()

	// delete votes
	new_votes := make([]int, 0)
	for v := range s.received_votes {
		if v != voter_id {
			new_votes = append(new_votes, v)
		}
	}
	s.received_votes = new_votes

	// send release
	rls_msg := RLSMessage(s.id, voter_id, s.clock, message_type)
	fmt.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
	logger.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
	s.servers[voter_id].ch <- rls_msg
}

// releaseAllVotes all votes after finish cs
func (s *Server) releaseAllVotes() {

	s.updateOwnClock()

	for i := 0; i < len(s.received_votes); i++ {
		voter_id := s.received_votes[i]
		rls_msg := RLSMessage(s.id, voter_id, s.clock, RLS)
		fmt.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
		logger.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
		s.servers[voter_id].ch <- rls_msg
	}
	s.received_votes = make([]int, 0)
}

// simulate execution of critical section
func (s *Server) executeCriticalSection() {
	s.Lock()
	s.executing_cs = true
	s.requesting = false
	s.Unlock()

	fmt.Printf("Server %d started cs execution.\n", s.id)
	logger.Printf("Server %d started cs execution.\n", s.id)

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

	s.Lock()
	s.executing_cs = false
	s.Unlock()

	fmt.Printf("Server %d has finished cs execution.\n", s.id)
	logger.Printf("Server %d has finished cs execution.\n", s.id)
}

// If already VOTEd for REQUEST(T’) and T’ is later than T
// then send RESCIND-VOTE to the machine requesting REQUEST(T’).
func (s *Server) onReceiveReq(msg Message) {
	req := &Request{
		clock:     msg.message.clock,
		requester: msg.message.requester,
	}
	if s.voted_req == nil {
		s.vote(msg.message.requester, msg.message.clock)
	} else if msg.message.clock < s.voted_req.clock || (msg.message.clock == s.voted_req.clock && msg.message.requester > s.voted_req.requester) {
		s.rescindVote(s.voted_req.requester)
		heap.Push(&s.pending_votes_queue, req)
	} else {
		heap.Push(&s.pending_votes_queue, req)
	}

}

func (s *Server) onReceiveVote(voter_id int) {
	if s.requesting {
		// check whether has received rescind previously
		for i := 0; i < len(s.received_rescinds); i++ {
			if s.received_rescinds[i] == voter_id {
				// remove from list
				s.received_rescinds[i] = s.received_rescinds[len(s.received_rescinds)-1]
				s.received_rescinds = s.received_rescinds[:len(s.received_rescinds)-1]
				// release and return
				s.releaseVote(voter_id, RESCIND_RLS)
				return
			}
		}

		s.Lock()
		s.received_votes = append(s.received_votes, voter_id)
		s.Unlock()

		if s.hasMajorityVotes() {
			logger.Printf("Server %d received_votes %d.\n", s.id, s.received_votes)

			s.executeCriticalSection()
			s.releaseAllVotes()
		}
	} else if s.executing_cs {
		s.received_votes = append(s.received_votes, voter_id)
	} else {
		s.releaseVote(voter_id, RLS)
	}
}

func (s *Server) onReceiveRls(releaser_id int, message_type string) {
	if message_type == RESCIND_RLS {
		if releaser_id == s.voted_req.requester {
			heap.Push(&s.pending_votes_queue, s.voted_req)
			if s.pending_votes_queue.Peek() != nil {
				request := heap.Pop(&s.pending_votes_queue).(*Request)
				s.vote(request.requester, request.clock)
			}
		}
	} else {
		if releaser_id == s.voted_req.requester {
			s.voted_req = nil
			if s.pending_votes_queue.Peek() != nil {
				request := heap.Pop(&s.pending_votes_queue).(*Request)
				s.vote(request.requester, request.clock)
			}
		}
	}
}

// If not already in the critical section, then send RELEASE-VOTE to the machine sending RESCIND-VOTE.
// Re-REQUEST for which the RESCIND-VOTE was sent
func (s *Server) onReceiveRcd(voter_id int) {
	if !s.executing_cs || !s.requesting {
		return
	}
	for i := 0; i < len(s.received_votes); i++ {
		if s.received_votes[i] == voter_id {
			s.releaseVote(voter_id, RESCIND_RLS)
			return
		}
	}

	// if hasn't replied vote, record first
	s.received_rescinds = append(s.received_rescinds, voter_id)
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
