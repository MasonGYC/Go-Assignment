package main

type Server struct {
	id    int
	ch    chan Message
	clock []int
}

type Item struct {
	value    interface{}
	priority int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

// pq := make(PriorityQueue, 0)

// 	// Push items with values and priorities into the queue
// 	items := []*Item{
// 		{value: "item1", priority: 3},
// 		{value: "item2", priority: 1},
// 		{value: "item3", priority: 2},
// 	}

// 	for _, item := range items {
// 		heap.Push(&pq, item)
// 	}

// 	// Pop items in priority order
// 	for pq.Len() > 0 {
// 		item := heap.Pop(&pq).(*Item)
// 		fmt.Printf("Item: %v, Priority: %d\n", item.value, item.priority)
// 	}
