package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
)

func main() {

	// log outputs for debugging purpose
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
	log.Printf("===============START===============")

	// define the number of servers
	num_servers := flag.Int("servers", 5, "number of servers")

	// define the number of concurrent requests to make
	// num_requests := flag.Int("requests", 2, "number of concurrent requests to make")
	flag.Parse()

	// initialize servers
	var servers = make([]*Server, *num_servers)
	for i := 0; i < *num_servers; i++ {
		servers[i] = NewServer(i, *num_servers, servers)
	}

	// start listening on all servers
	for i := 0; i < *num_servers; i++ {
		go servers[i].listen()
	}

	random_server_idx := rand.Intn(len(servers))
	go servers[random_server_idx].request()

	// TODO: to be filled

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
