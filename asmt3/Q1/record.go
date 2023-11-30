package main

import "container/heap"

// access
const (
	NIL = 0
	R   = 1
	RW  = 2
)

type ManagerRecord struct {
	page_num int
	copySet  []int
	owner    int
	writing  bool
	queue    PriorityQueue
}

func NewManagerRecord(pageNum int, copySet []int, owner int) ManagerRecord {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return ManagerRecord{
		page_num: pageNum,
		copySet:  copySet,
		owner:    owner,
		writing:  false,
		queue:    pq,
	}
}

type ServerRecord struct {
	page_num int
	access   int
	page     Page
}

func NewServerRecord(pageNum int, access int, page Page) ServerRecord {
	return ServerRecord{
		page_num: pageNum,
		access:   access,
		page:     page,
	}
}
