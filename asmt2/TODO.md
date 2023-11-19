# assumptions:
priority has 2 metrics: 
1. clock (smaller clock -> higher priority)
2. node id (higher id -> higher priority)

every server can only make one request at one time. means if doesn't receive all replies, it won't start making new request. 

clocks are updated before the execution of a certain action

go run LSPQ.go logger.go message.go PriorityQueue.go server.go
go run LSPQ_RA.go logger.go message.go PriorityQueue.go server.go
go run voting.go logger.go message.go queue.go server.go clock.go

nanoseconds
