package main

// 2. Use Lamport’s logical clock to determine a total order of all the messages received at
// all the registered clients. Subsequently, present (i.e., print) this order for all registered
// clients to know the order in which the messages should be read. (10 points)

// 3. Use Vector clock to redo the assignment. Implement the detection of causality violation
// and print any such detected causality violation. (10 points)

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	// Start a goroutine toreceive server forward message
	go func() {
		for {

			// receive data
			m := <-c.ch
			fmt.Println("Client "+strconv.Itoa(c.id)+" received:", m)
			log.Println("Client "+strconv.Itoa(c.id)+" received:", m)

			// update clock
			mutex.Lock()
			c.clock = max(m.clock, c.clock) + 1
			mutex.Unlock()

			fmt.Println("Client "+strconv.Itoa(c.id)+"'s clock:", c.clock)
			log.Println("Client "+strconv.Itoa(c.id)+"'s clock:", c.clock)
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

			fmt.Println("Client "+strconv.Itoa(c.id)+"'s clock:", c.clock)
			log.Println("Client "+strconv.Itoa(c.id)+"'s clock:", c.clock)

			fmt.Println("Client "+strconv.Itoa(c.id)+" sending:", send)
			log.Println("Client "+strconv.Itoa(c.id)+" sending:", send)

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
		fmt.Println("Server received:", msg)
		log.Println("Server received:", msg)

		// update clock
		mutex.Lock()
		s.clock = max(msg.clock, s.clock) + 1
		mutex.Unlock()

		fmt.Println("Server's clock: ", s.clock)
		log.Println("Server's clock: ", s.clock)

		// check forward or drop
		flag := rand.Intn(2)

		if flag == 0 {
			// forward to clients except sender
			msg.clock = s.clock
			for j := 0; j < s.num_clients; j++ {
				// check if not sender, then forward
				if j != msg.sender_id {
					fmt.Println("Server forward to Client "+strconv.Itoa(j), msg)
					log.Println("Server forward to Client "+strconv.Itoa(j), msg)

					s.clients[j].ch <- msg
				}

			}
		} else {
			fmt.Println("Drop the message", msg)
			log.Println("Drop the message", msg)
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
	log.Println("===============START===============")

	// define the number of clients
	const num_clients int = 15

	// init server and clients
	var server = Server{
		num_clients: num_clients,
		clients:     make(map[int]Client),
		ch:          make(chan Message),
		clock:       0,
	}

	for i := 0; i < num_clients; i++ {
		server.clients[i] = newClient(i, make(chan Message), 0)
	}

	// execute client and server go routines
	for i := 0; i < num_clients; i++ {
		go clientSend(server.clients[i], server.ch)
	}
	go serverRecv(server)

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
