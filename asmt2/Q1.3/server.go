package main

import (
	"fmt"
	"sync"
	"time"
)

type Server struct {
	id                  int
	ch                  chan Message
	clock               []int
	servers             []*Server // all servers in the network
	voted_id            int       // if -1, not voted; else voted to who
	voted_req_clock     []int     // the clock of the voted request
	requested           bool      // true if requested
	received_votes      []int     // id of voters; len() = num of votes received
	pending_votes_queue MsgQueue  // votes to be sent
	executing_cs        bool      // is executing cs
	sync.Mutex
}

func NewServer(id int, num_servers int, servers []*Server) *Server {
	return &Server{
		id:                  id,
		ch:                  make(chan Message),
		clock:               make([]int, num_servers),
		servers:             servers,
		voted_id:            -1,
		requested:           false,
		received_votes:      make([]int, 0),
		pending_votes_queue: *NewQueue(),
		executing_cs:        false,
	}
}

func (s *Server) hasMajorityVotes() bool {
	received_votes := len(s.received_votes)
	num_of_servers := len(s.servers)
	if num_of_servers%2 == 1 {
		num_of_servers = num_of_servers - 1
	}
	return received_votes > num_of_servers/2
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

				go s.onReceiveRls(msg.sender_id)
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

	// self-vote if not voted
	s.Lock()
	s.requested = true
	if s.voted_id == -1 {
		s.received_votes = append(s.received_votes, s.id)
	}
	s.Unlock()

	// broadcast request
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

// vote
func (s *Server) vote(requester_id int, request_clock []int) {

	s.updateOwnClock()

	s.Lock()
	s.voted_id = requester_id
	s.voted_req_clock = request_clock
	s.Unlock()

	vote_msg := VOTEMessage(s.id, requester_id, s.clock)
	s.servers[requester_id].ch <- vote_msg

	fmt.Printf("Server %d vote to %d at %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d vote to %d at %d.\n", s.id, requester_id, s.clock)
}

// rescind vote
func (s *Server) rescind_vote(requester_id int) {

	s.updateOwnClock()

	rcd_msg := RESCINDMessage(s.id, requester_id, s.clock)
	s.servers[requester_id].ch <- rcd_msg

	fmt.Printf("Server %d rescind vote to %d at %d.\n", s.id, requester_id, s.clock)
	logger.Printf("Server %d rescind vote to %d at %d.\n", s.id, requester_id, s.clock)
}

// release vote from voter
func (s *Server) releaseVote(voter_id int) {

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
	rls_msg := RLSMessage(s.id, voter_id, s.clock)
	s.servers[voter_id].ch <- rls_msg

	fmt.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
	logger.Printf("Server %d release to %d at %d.\n", s.id, voter_id, s.clock)
}

// releaseAllVotes all votes after finish cs
func (s *Server) releaseAllVotes() {

	s.updateOwnClock()

	servers := s.servers
	for i := 0; i < len(servers); i++ {
		if servers[i].id != s.id {
			rls_msg := RLSMessage(s.id, servers[i].id, s.clock)
			servers[i].ch <- rls_msg

			fmt.Printf("Server %d release to %d at %d.\n", s.id, servers[i].id, s.clock)
			logger.Printf("Server %d release to %d at %d.\n", s.id, servers[i].id, s.clock)
		}
	}

	// reset request records
	s.Lock()
	s.requested = false
	s.received_votes = make([]int, 0)
	s.Unlock()
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
	if s.voted_id == -1 {
		s.vote(msg.sender_id, msg.clock)
	} else {
		// If already VOTEd for REQUEST(T’) and T’ is later than T
		// then send RESCIND-VOTE to the machine requesting REQUEST(T’).
		if clockLess(msg.clock, s.voted_req_clock) {
			s.rescind_vote(msg.sender_id)
		} else {
			// put in queue
			s.pending_votes_queue.Push(msg)
		}
	}
}

func (s *Server) onReceiveVote(voter_id int) {
	s.Lock()
	s.received_votes = append(s.received_votes, voter_id)
	s.Unlock()
	if s.hasMajorityVotes() && !s.executing_cs {
		s.executing_cs = true
		s.executeCriticalSection()
		s.executing_cs = false
		s.releaseAllVotes()
	}
}

func (s *Server) onReceiveRls(releaser_id int) {
	s.Lock()
	s.voted_id = -1
	s.Unlock()

	// check for pending votes
	if !s.pending_votes_queue.IsEmpty() {
		msg := s.pending_votes_queue.Pop()
		s.vote(msg.sender_id, msg.clock)
	}

}

// If not already in the critical section, then send RELEASE-VOTE to the machine sending RESCIND-VOTE.
// Re-REQUEST for which the RESCIND-VOTE was sent
func (s *Server) onReceiveRcd(voter_id int) {
	if !s.executing_cs {
		s.releaseVote(voter_id)
	}
}

// update clock
func (s *Server) updateClock(msgClock []int) {
	s.Lock()
	for i := 0; i < len(s.clock); i++ {
		s.clock[i] = max(msgClock[i], s.clock[i])
	}
	s.clock[s.id] = s.clock[s.id] + 1
	s.Unlock()

	// fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}

// update clock
func (s *Server) updateOwnClock() {
	s.Lock()
	s.clock[s.id] = s.clock[s.id] + 1
	s.Unlock()

	// fmt.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
	// logger.Printf("Server %d's clock updated: %d.\n", s.id, s.clock)
}
