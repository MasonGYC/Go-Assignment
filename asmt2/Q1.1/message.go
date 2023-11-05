package main

import "fmt"

const (
	// message types
	REQ string = "REQUEST"
	RES string = "REPLY"
	RLS string = "RELEASE"
)

type Message struct {
	sender_id    int
	receiver_id  int
	message      string //content
	message_type string
	clock        []int // vector clock
}

// request to access cs
func REQMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Request from server %d to server %d to enter critical section.", sender_id, receiver_id),
		message_type: REQ,
		clock:        clock,
	}
}

// reply to grant permission
func RESMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Reply from server %d to server %d to grant permission.", sender_id, receiver_id),
		message_type: RES,
		clock:        clock,
	}
}

// relase cs
func RLSMessage(sender_id int, receiver_id int, clock []int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      fmt.Sprintf("Server %d has relased cs. Server %d is informed.", sender_id, receiver_id),
		message_type: RLS,
		clock:        clock,
	}
}
