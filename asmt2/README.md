# assumptions:
priority has 2 metrics: 
1. vetvor clock (smaller clock higher priority)
2. node id (higher id higher priority)

every server can only make one request at one time. means if doesn't receive all replies, it won't start making new request. 

clocks are updated before the execution of a certain action

`reply_counter` in `Request`
- for requester: `len(servers) - 1`
- for receiver: `1` if not replied, `0` if replied

unique identifier for each request: requester_id and clock
./run.exe -servers

# TODO:
1. add clock update
2. multiple concurrent requests (hold reply)
3. start with random clock
4. terminate automatically when all finished (waitgroup?)


go run LSPQ.go clock.go logger.go message.go PriorityQueue.go server.go
go run voting.go clock.go logger.go message.go queue.go server.go