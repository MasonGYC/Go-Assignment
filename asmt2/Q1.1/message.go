package main

import (
	"fmt"
)

const (
	// message types
	REQ string = "REQUEST"
	RES string = "REPLY"
	RLS string = "RELEASE"
)

type Message struct {
	sender_id    int
	receiver_id  int
	message      any //content
	message_type string
	clock        []int // vector clock
}

// request to access cs
func REQMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Request from node %d to node %d to enter critical section.", sender_id, receiver_id),
		message_type: REQ,
		clock:        clock,
	}
}

// reply to grant permission
func RESMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Reply from node %d to node %d to grant permission.", sender_id, receiver_id),
		message_type: RES,
		clock:        clock,
	}
}

// relase cs
func RLSMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Node %d has relased cs. Node %d is informed.", sender_id, receiver_id),
		message_type: RLS,
		clock:        clock,
	}
}
