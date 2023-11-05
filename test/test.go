// This example demonstrates a priority queue built using the heap interface.
package main

import (
	"fmt"
)

type RequestPriorityQueue []*int

func main() {
	pq := make(RequestPriorityQueue, 0)
	fmt.Printf("%d", pq[0])
}
