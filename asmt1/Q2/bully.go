// // 1. Implement the above protocol of joint synchronization and election via GO (20 points)
// // The coordinator initiates the synchronization process by sending message to all other machines.
// // Upon receiving the message from the coordinator, each machine updates its local version of the data structure with the coordinator’s version.
// // a new coordinator is chosen by the Bully algorithm. You can assume a fixed timeout to simulate the behaviour of detecting a fault.
package main

// version working w/o bully. can send sync msg

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const DATA int = 0
const SYNC int = 1

type Data struct {
}

type Replica struct {
	id      int
	ch      chan Message
	data    Data
	timeout time.Duration // assume if not receiving message after timeout, then coordinator failed
}

type Coordinator struct {
	id            int
	num_replicas  int
	replicas      map[int]Replica
	ch            chan Message
	data          Data
	sync_interval time.Duration //send sync msg every sync_interval
}

type Message struct {
	sender_id    int
	message_type int
	message      any // i-th msg
}

func bully() {

}

func replica(rep Replica, ch chan Message) {

	// Start a goroutine to receive coordinator message
	go func() {
		for {

			// receive data
			select {
			case msg := <-rep.ch:

				fmt.Println("Replica "+strconv.Itoa(rep.id)+" received:", msg)
				log.Println("Replica "+strconv.Itoa(rep.id)+" received:", msg)

				switch msg.message_type {
				case SYNC:
					coor_data := msg.message.(Data) // SYNC come with Data in message field
					rep.data = coor_data            // update local data
				default:

				}
			case <-time.After(rep.timeout):

				fmt.Println("Coordinator failed")
				log.Println("Coordinator failed")

				go bully() //select new coordinator using bully algo
			}

		}
	}()
}

func coordinator(coor Coordinator) {

	// initiates the synchronization process
	go func() {
		for {

			// create sync message
			var msg = Message{
				sender_id:    coor.id,
				message_type: SYNC,
				message:      coor.data,
			}

			// send sync message to all other machines.
			for i := 0; i < coor.num_replicas; i++ {
				coor.replicas[i].ch <- msg
			}

			// sleep for sync_interval
			time.Sleep(coor.sync_interval)
		}

	}()

}

func main() {

	// log outputs for debugging purpose
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Println("===============START===============")

	// define the number of clients
	const num_replicas int = 15

	// init clients and channels
	var coor = Coordinator{
		id:            0, // first one
		num_replicas:  num_replicas,
		replicas:      make(map[int]Replica),
		ch:            make(chan Message),
		sync_interval: time.Second,
	}

	for i := 0; i < num_replicas; i++ {
		coor.replicas[i] = Replica{
			id:      i + 1,
			ch:      make(chan Message),
			timeout: 2 * time.Second,
		}
	}

	// execute client and server go routines
	for i := 0; i < num_replicas; i++ {
		go replica(coor.replicas[i], coor.ch)
	}
	go coordinator(coor)

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}

// //TODO: a function to simulate coordinator fails
