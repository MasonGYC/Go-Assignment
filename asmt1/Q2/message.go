package main

import (
	"fmt"
	"time"
)

type Message struct {
	sender_id    int
	message_type string
	message      any
	timestamp    time.Time
}

// refuse self-elect request
func ACKMessage(request_id int, refuse_id int) Message {
	return Message{
		sender_id:    refuse_id,
		message_type: ACK,
		message:      fmt.Sprintf("Request from node %d is refused by node %d.", request_id, refuse_id),
		timestamp:    time.Now(),
	}
}

// self elect message
func ELECTMessage(sender_id int, new_coor_id int) Message {
	return Message{
		sender_id:    sender_id,
		message_type: ELECT,
		message:      fmt.Sprintf("Elect %d as new coordinator.", new_coor_id),
		timestamp:    time.Now(),
	}
}

// coordinator send this message to sync data with worker
func SYNCMessage(sender_id int, data Data) Message {
	return Message{
		sender_id:    sender_id,
		message_type: SYNC,
		message:      data,
		timestamp:    time.Now(),
	}
}

// newly elected coordinator broadcast its win
func VICTORYMessage(sender_id int, winner_id int) Message {
	return Message{
		sender_id:    sender_id,
		message_type: VICTORY,
		message:      fmt.Sprintf("Node %d is the new coordinator", winner_id),
		timestamp:    time.Now(),
	}
}

// Stop a goroutine
func STOPMessage(sender_id int) Message {
	return Message{
		sender_id:    sender_id,
		message_type: STOP,
		message:      "Stop the goroutine",
		timestamp:    time.Now(),
	}
}
