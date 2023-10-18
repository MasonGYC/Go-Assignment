package main

import (
	"fmt"
	"time"
)

type Message struct {
	sender_id    int
	message_type string
	message      any // i-th msg
	timestamp    time.Time
}

func ACKMessage(request_id int, refuse_id int) Message {
	return Message{
		sender_id:    refuse_id,
		message_type: ACK,
		message:      fmt.Sprintf("Request from node %d is refused by node %d.", request_id, refuse_id),
		timestamp:    time.Now(),
	}
}

// only self elect is possible can broadcase, so sender_id == new_coor_id
func ELECTMessage(sender_id int, new_coor_id int) Message {
	return Message{
		sender_id:    sender_id,
		message_type: ELECT,
		message:      fmt.Sprintf("Elect %d as new coordinator.", new_coor_id),
		timestamp:    time.Now(),
	}
}

func SYNCMessage(sender_id int, data Data) Message {
	return Message{
		sender_id:    sender_id,
		message_type: SYNC,
		message:      data,
		timestamp:    time.Now(),
	}
}

// only winner can broadcase, so sender_id == winner_id
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
