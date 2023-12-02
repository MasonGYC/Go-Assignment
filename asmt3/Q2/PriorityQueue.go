// Adapted from https://pkg.go.dev/container/heap
package main

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Message

func (pq PriorityQueue) Len() int { return len(pq) }

// lower clock or higher sender_id (if concurrent) has higher priority.
func (pq PriorityQueue) Less(i, j int) bool {
	if pq[i].clock < pq[j].clock || pq[j].clock < pq[i].clock {
		return pq[i].clock < pq[j].clock
	} else {
		return pq[i].sender_id > pq[j].sender_id
	}
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	// pq[i].index = i
	// pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	// n := len(*pq)
	item := x.(*Message)
	// item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	// item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// // update modifies the priority and value of an Item in the queue.
// func (pq *PriorityQueue) Update(item *Message, value string, clock int) {
// 	item.value = value
// 	item.clock = clock
// 	heap.Fix(pq, item.index)
// }

// peek the head of the queue
func (pq PriorityQueue) Peek() *Message {
	if len(pq) > 0 {
		return pq[0]
	}
	return nil
}
