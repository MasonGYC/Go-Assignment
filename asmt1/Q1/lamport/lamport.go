package main

import (
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
	clock int
}

func newClient(id int, ch chan Message, clock int) Client {
	return Client{
		id:    id,
		ch:    ch,
		clock: clock,
	}
}

type Server struct {
	num_clients int
	clients     map[int]Client
	ch          chan Message
	clock       int
}

type Message struct {
	sender_id int
	message   int // i-th msg
	clock     int //lamport's logical clock
}

func newMessage(sender_id int, message int, clock int) Message {
	return Message{
		sender_id: sender_id,
		message:   message,
		clock:     clock,
	}
}

func clientSend(c Client, sch chan Message) {

	var mutex sync.Mutex

	// Start a goroutine to receive server forward message
	go func() {
		for {

			// receive data
			m := <-c.ch
			fmt.Printf("Client %d received Server's message. Sent from Client %d at clock %d.\n", c.id, m.sender_id, m.clock)
			log.Printf("Client %d received Server's message. Sent from Client %d at clock %d.\n", c.id, m.sender_id, m.clock)

			// update clock
			mutex.Lock()
			c.clock = max(m.clock, c.clock) + 1
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
			c.clock++
			mutex.Unlock()

			// send a new message
			send := newMessage(c.id, i, c.clock)

			fmt.Printf("Client %d's clock:  %d.\n", c.id, c.clock)
			log.Printf("Client %d's clock:  %d.\n", c.id, c.clock)

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

		// update clock
		mutex.Lock()
		s.clock = max(msg.clock, s.clock) + 1
		mutex.Unlock()

		fmt.Printf("Server's clock: %d.\n", s.clock)
		log.Printf("Server's clock: %d.\n", s.clock)

		// check forward or drop, each has 50% chance
		flag := rand.Intn(2)

		if flag == 0 {
			// forward to clients except sender
			msg.clock = s.clock
			for j := 0; j < s.num_clients; j++ {
				// check if not sender, then forward
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

	// define the number of clients, pass by flag
	num_clients := flag.Int("clients", 15, "number of clients")
	flag.Parse()

	// init server and clients
	var server = Server{
		num_clients: *num_clients,
		clients:     make(map[int]Client),
		ch:          make(chan Message),
		clock:       0,
	}

	for i := 0; i < server.num_clients; i++ {
		server.clients[i] = newClient(i, make(chan Message), 0)
	}

	// execute client goroutines
	for i := 0; i < server.num_clients; i++ {
		go clientSend(server.clients[i], server.ch)
	}
	// execute server goroutines
	go serverRecv(server)

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
