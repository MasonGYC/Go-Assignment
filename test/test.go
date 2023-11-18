// Adapted from https://pkg.go.dev/container/heap
package main

import (
	"container/heap"
	"fmt"
)

// An Request is something we manage in a priority queue.
type Request struct {
	value     string // the sender's i-th request
	clock     int
	requester int // the id of the requester
	index     int // The index of the item in the heap.
}

type RequestIdentifier struct {
	clock     int
	requester int // the id of the requester
}

func (r1 *Request) isEqual(r2 Request) bool {
	if r1.value == r2.value && r1.clock == r2.clock && r1.requester == r2.requester {
		return true
	} else {
		return false
	}
}

// A RequestPriorityQueue implements heap.Interface and holds Items.
type RequestPriorityQueue []*Request

func (pq RequestPriorityQueue) Len() int { return len(pq) }

// lower clock or higher sender_id (if concurrent) has higher priority.
func (pq RequestPriorityQueue) Less(i, j int) bool {
	if pq[i].clock < pq[j].clock || pq[j].clock < pq[i].clock {
		return pq[i].clock < pq[j].clock
	} else {
		return pq[i].requester > pq[j].requester
	}
}

func (pq RequestPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *RequestPriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Request)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *RequestPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *RequestPriorityQueue) Update(item *Request, value string, clock int) {
	item.value = value
	item.clock = clock
	heap.Fix(pq, item.index)
}

// peek the head of the queue
func (pq RequestPriorityQueue) Peek() *Request {
	if len(pq) > 0 {
		return pq[0]
	}
	return nil
}

// This example creates a PriorityQueue with some items, adds and manipulates an item,
// and then removes the items in priority order.
func main() {
	// Some items and their priorities.
	items := map[string]int{
		"s0": 3, "s1": 1, "s2": 1,
	}

	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(RequestPriorityQueue, len(items))
	i := 0
	for value, clock := range items {
		pq[i] = &Request{
			value:     value,
			clock:     clock,
			requester: i,
			index:     i,
		}
		i++
	}
	heap.Init(&pq)

	// Insert a new item and then modify its priority.
	item := &Request{
		value:     "s3",
		clock:     1,
		requester: 3,
	}
	heap.Push(&pq, item)

	// Take the items out; they arrive in decreasing priority order.
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Request)
		fmt.Printf("%.2d:%s ", item.clock, item.value)
	}
}
