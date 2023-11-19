package main

const (
	// message types
	REQ         string = "REQUEST-VOTE"
	VOTE        string = "VOTE"
	RLS         string = "RELEASE-VOTE"
	RESCIND     string = "RESCIND-VOTE"
	RESCIND_RLS string = "RESCIND_RLS"
)

type Message struct {
	sender_id    int
	receiver_id  int
	message      RequestIdentifier
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

// vote for the request
func VOTEMessage(sender_id int, receiver_id int, request_clock int, clock int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      RequestIdentifier{clock: request_clock, requester: receiver_id},
		message_type: VOTE,
		clock:        clock,
	}
}

// release vote
func RLSMessage(sender_id int, receiver_id int, clock int, message_type string) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      RequestIdentifier{clock: clock, requester: receiver_id},
		message_type: message_type,
		clock:        clock,
	}
}

// rescind vote
func RESCINDMessage(sender_id int, receiver_id int, request_clock int, clock int) Message {
	return Message{
		sender_id:    sender_id,
		receiver_id:  receiver_id,
		message:      RequestIdentifier{clock: request_clock, requester: receiver_id},
		message_type: RESCIND,
		clock:        clock,
	}
}
