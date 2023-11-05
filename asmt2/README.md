# assumptions:
priority has 2 metrics: 
1. vetvor clock (smaller clock higher priority)
2. node id (higher id higher priority)

every server can only make one request at one time. means if doesn't receive all replies, it won't start making new request. 

clocks are updated before the execution of a certain action

./run.exe -servers

# TODO:
1. add clock update
2. multiple concurrent requests (hold reply)


go run LSPQ.go clock.go logger.go message.go PriorityQueue.go server.go