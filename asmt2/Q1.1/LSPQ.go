package main

import (
	"fmt"
	"log"
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

	// TODO: to be filled

	var input string
	// wait for the input, as otherwise, the program will not wait
	fmt.Scanln(&input)
}
