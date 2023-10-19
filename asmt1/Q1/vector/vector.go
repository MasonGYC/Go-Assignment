package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Client struct {
	id    int
	ch    chan Message
	clock []int
}

func newClient(id int, ch chan Message, clock []int) Client {
	return Client{
		id:    id,
		ch:    ch,
		clock: clock,
	}
}

type Server struct {
	num_clients int
	id          int
	clients     map[int]Client
	ch          chan Message
	clock       []int
}

type Message struct {
	sender_id int
	message   int   // i-th msg
	clock     []int //vector clock
}

func newMessage(sender_id int, message int, clock []int) Message {
	return Message{
		sender_id: sender_id,
		message:   message,
		clock:     clock,
	}
}

// check whether c1 > c2
func clockIsGreaterThan(c1 []int, c2 []int) (bool, error) {
	// check comparability
	if len(c1) == len(c2) {
		// loop thru to see if all >=
		for i := 0; i < len(c1); i++ {
			if c1[i] < c2[i] {
				// if there is less than, false
				return false, nil
			}
		}
		// loop thru to see if there exist >
		for i := 0; i < len(c1); i++ {
			if c1[i] > c2[i] {
				// since all >=, if one >, then true
				return true, nil
			}
		}
		// since all >=, there's no >, then false
		return false, nil
	} else {
		// if not comparable, error
		return false, errors.New("2 clocks have different number of entries")
	}

}

func clientSend(c Client, sch chan Message) {

	var mutex sync.Mutex

	// Start a goroutine to receive server forward message
	go func() {
		for {

			// receive data
			msg := <-c.ch

			fmt.Printf("Client %d received Server's message. Sent from Client %d at clock %d.\n", c.id, msg.sender_id, msg.clock)
			log.Printf("Client %d received Server's message. Sent from Client %d at clock %d.\n", c.id, msg.sender_id, msg.clock)

			// check causality violation
			// If the local vector clock of the receiving machine is more than the vector clock Of the message, then a potential causality violation is detected.
			violation, _ := clockIsGreaterThan(c.clock, msg.clock)
			if violation {
				fmt.Printf("Potential violation detected!")
				fmt.Printf("Local clock: %d.\n", c.clock)
				fmt.Printf("Message clock: %d.\n", msg.clock)
				log.Printf("Potential violation detected!")
				log.Printf("Local clock: %d.\n", c.clock)
				log.Printf("Message clock: %d.\n", msg.clock)
			}

			// update clock
			mutex.Lock()
			for i := 0; i < len(c.clock); i++ {
				c.clock[i] = max(msg.clock[i], c.clock[i])
			}
			c.clock[c.id] = c.clock[c.id] + 1
			mutex.Unlock()

			fmt.Printf("Client %d's clock: %d.\n", c.id, c.clock)
			log.Printf("Client %d's clock: %d.\n", c.id, c.clock)
		}
	}()

	// Start a goroutine to send message to the server
	go func() {

		for i := 0; ; i++ {

			// update clock
			mutex.Lock()
			c.clock[c.id] = c.clock[c.id] + 1
			mutex.Unlock()

			// send a new message
			send := newMessage(c.id, i, c.clock)

			fmt.Printf("Client %d's clock: %d.\n", c.id, c.clock)
			log.Printf("Client %d's clock: %d.\n", c.id, c.clock)

			fmt.Printf("Client %d sending the %d-th message at clock %d.\n", c.id, send.message, send.clock)
			log.Printf("Client %d sending the %d-th message at clock %d.\n", c.id, send.message, send.clock)

			sch <- send

			// sleep for a while
			amt := time.Duration(rand.Intn(5000))
			time.Sleep(time.Millisecond * amt)
		}
	}()
}

func serverRecv(s Server) {

	var mutex sync.Mutex

	for {
		// receive messages from all channels and print it
		msg := <-s.ch
		fmt.Printf("Server received Client %d's %d-th message of clock %d.\n", msg.sender_id, msg.message, msg.clock)
		log.Printf("Server received Client %d's %d-th message of clock %d.\n", msg.sender_id, msg.message, msg.clock)

		// check causality violation
		// If the local vector clock of the receiving machine is more than the vector clock Of the message, then a potential causality violation is detected.
		violation, _ := clockIsGreaterThan(s.clock, msg.clock)
		if violation {
			fmt.Printf("Potential violation detected!")
			fmt.Printf("Local clock: %d.\n", s.clock)
			fmt.Printf("Message clock: %d.\n", msg.clock)
			log.Printf("Potential violation detected!")
			log.Printf("Local clock: %d.\n", s.clock)
			log.Printf("Message clock: %d.\n", msg.clock)
		}

		// update clock
		mutex.Lock()
		for i := 0; i < len(s.clock); i++ {
			s.clock[i] = max(msg.clock[i], s.clock[i])
		}
		s.clock[s.id] = s.clock[s.id] + 1
		mutex.Unlock()

		fmt.Printf("Server's clock: %d.\n", s.clock)
		log.Printf("Server's clock: %d.\n", s.clock)

		// check forward or drop
		flag := rand.Intn(2)

		if flag == 0 {
			// forward to clients except sender
			msg.clock = s.clock
			for j := 0; j < s.num_clients; j++ {
				if j != msg.sender_id {
					fmt.Printf("Server forward to Client %d Client %d's %d-th message.\n", j, msg.sender_id, msg.message)
					log.Printf("Server forward to Client %d Client %d's %d-th message.\n", j, msg.sender_id, msg.message)

					s.clients[j].ch <- msg
				}
			}
		} else {
			fmt.Printf("Drop Client %d's %d-th message of clock %d.\n", msg.sender_id, msg.message, msg.clock)
			log.Printf("Drop Client %d's %d-th message of clock %d.\n", msg.sender_id, msg.message, msg.clock)
		}

	}
}

func main() {

	// log outputs for debugging purpose
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Printf("===============START===============")

	// define the number of clients
	num_clients := flag.Int("clients", 15, "number of clients")
	flag.Parse()

	// init server and clients
	var server = Server{
		num_clients: *num_clients,
		id:          *num_clients, // last one
		clients:     make(map[int]Client),
		ch:          make(chan Message),
		clock:       make([]int, *num_clients+1), //plus server
	}

	for i := 0; i < server.num_clients; i++ {
		server.clients[i] = newClient(i, make(chan Message), make([]int, server.num_clients+1))
	}

	// execute client and server go routines
	for i := 0; i < server.num_clients; i++ {
		go clientSend(server.clients[i], server.ch)
	}
	go serverRecv(server)

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
