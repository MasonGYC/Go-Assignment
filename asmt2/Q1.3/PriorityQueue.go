// Adapted from https://pkg.go.dev/container/heap
package main

import (
	"container/heap"
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
