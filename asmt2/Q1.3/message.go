package main

import "fmt"

const (
	// message types
	REQ     string = "REQUEST-VOTE"
	VOTE    string = "VOTE"
	RLS     string = "RELEASE-VOTE"
	RESCIND string = "RESCIND-VOTE"
)

type Message struct {
	sender_id    int
	receiver_id  int
	message      any //TODO: ?
	message_type string
	clock        []int // vector clock
}

// request to access cs
func REQMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Server %d reuqets to server %d to enter cs.\n", sender_id, receiver_id),
		message_type: REQ,
		clock:        clock,
	}
}

// vote for the request
func VOTEMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Server %d votes for server %d.\n", sender_id, receiver_id),
		message_type: VOTE,
		clock:        clock,
	}
}

// release vote
func RLSMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Server %d released vote to server %d.\n", sender_id, receiver_id),
		message_type: RLS,
		clock:        clock,
	}
}

// rescind vote
func RESCINDMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Server %d rescind vote to server %d.\n", sender_id, receiver_id),
		message_type: RESCIND,
		clock:        clock,
	}
}
