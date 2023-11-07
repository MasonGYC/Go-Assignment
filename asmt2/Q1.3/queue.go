package main

type MsgQueue struct {
	queue []Message
}

func NewQueue() *MsgQueue {
	return &MsgQueue{
		queue: make([]Message, 0),
	}
}

func (q *MsgQueue) Push(m Message) {
	q.queue = append(q.queue, m)
}

func (q *MsgQueue) Peek() Message {
	return q.queue[0]
}

func (q *MsgQueue) Pop() Message {
	top := q.queue[0]
	q.queue = q.queue[1:]
	return top
}

func (q *MsgQueue) IsEmpty() bool {
	return len(q.queue) == 0
}
