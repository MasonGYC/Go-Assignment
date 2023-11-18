package main

const (
	// message types
	REQ string = "REQUEST"
	REP string = "REPLY"
)

type Message struct {
	sender_id    int
	receiver_id  int
	message      RequestIdentifier //requester_id: request_clock
	message_type string
	clock        int // vector clock
}

// request to access cs
func REQMessage(sender_id int, receiver_id int, request_clock int, clock int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      RequestIdentifier{clock: request_clock, requester: sender_id},
		message_type: REQ,
		clock:        clock,
	}
}

// reply to grant permission
func REPMessage(sender_id int, receiver_id int, request_clock int, clock int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      RequestIdentifier{clock: request_clock, requester: receiver_id},
		message_type: REP,
		clock:        clock,
	}
}
