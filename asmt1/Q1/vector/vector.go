// 3. Use Vector clock to redo the assignment. Implement the detection of causality violation
// and print any such detected causality violation. (10 points)

package main

import (
	"errors"
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

			fmt.Println("Client "+strconv.Itoa(c.id)+" received:", msg)
			log.Println("Client "+strconv.Itoa(c.id)+" received:", msg)

			// check causality violation
			// If the local vector clock of the receiving machine is more than the vector clock Of the message, then a potential causality violation is detected.
			violation, _ := clockIsGreaterThan(c.clock, msg.clock)
			if violation {
				fmt.Println("Potential violation detected!")
				fmt.Println("Local clock: ", c.clock)
				fmt.Println("Message clock: ", msg.clock)
				log.Println("Potential violation detected!")
				log.Println("Local clock: ", c.clock)
				log.Println("Message clock: ", msg.clock)
			}

			// update clock
			mutex.Lock()
			for i := 0; i < len(c.clock); i++ {
				c.clock[i] = max(msg.clock[i], c.clock[i])
			}
			c.clock[c.id] = c.clock[c.id] + 1
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
			c.clock[c.id] = c.clock[c.id] + 1
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

		// check causality violation
		// If the local vector clock of the receiving machine is more than the vector clock Of the message, then a potential causality violation is detected.
		violation, _ := clockIsGreaterThan(s.clock, msg.clock)
		if violation {
			fmt.Println("Potential violation detected!")
			fmt.Println("Local clock: ", s.clock)
			fmt.Println("Message clock: ", msg.clock)
			log.Println("Potential violation detected!")
			log.Println("Local clock: ", s.clock)
			log.Println("Message clock: ", msg.clock)
		}

		// update clock
		mutex.Lock()
		for i := 0; i < len(s.clock); i++ {
			s.clock[i] = max(msg.clock[i], s.clock[i])
		}
		s.clock[s.id] = s.clock[s.id] + 1
		mutex.Unlock()

		fmt.Println("Server's clock: ", s.clock)
		log.Println("Server's clock: ", s.clock)

		// check forward or drop
		flag := rand.Intn(2)

		if flag == 0 {
			// forward to clients except sender
			msg.clock = s.clock
			for j := 0; j < s.num_clients; j++ {
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

	// init clients and channels
	var server = Server{
		num_clients: num_clients,
		id:          num_clients, // last one
		clients:     make(map[int]Client),
		ch:          make(chan Message),
		clock:       make([]int, num_clients+1), //plus server
	}

	for i := 0; i < num_clients; i++ {
		server.clients[i] = newClient(i, make(chan Message), make([]int, num_clients+1))
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
