package main

const (
	// message types
	SYNC    string = "SYNC"
	ELECT   string = "ELECT"
	ACK     string = "ACK"
	VICTORY string = "VICTORY"
	STOP    string = "STOP"

	// node status
	COORDINATOR string = "COORDINATOR"
	WORKER      string = "WORKER"

	// node state
	// only following states possile: RN,CN,RE,CD,RD,RB
	NORMAL        = "NORMAL"
	SELF_ELECTING = "SELF_ELECTING"
	BROADCATING   = "BROADCATING"
	DOWN          = "DOWN"
)

type Data struct {
}
