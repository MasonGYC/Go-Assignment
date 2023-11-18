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
4. terminate automatically when all finished (waitgroup?)
5. compare performance


go run LSPQ.go logger.go message.go PriorityQueue.go server.go
go run LSPQ_RA.go logger.go message.go PriorityQueue.go server.go
go run voting.go logger.go message.go queue.go server.go clock.go

nanoseconds


cd ../Q1.2;
LSPQ_RA.exe -servers=10 -requests=1
LSPQ_RA.exe -servers=10 -requests=2
LSPQ_RA.exe -servers=10 -requests=3
LSPQ_RA.exe -servers=10 -requests=4
LSPQ_RA.exe -servers=10 -requests=5
LSPQ_RA.exe -servers=10 -requests=6
LSPQ_RA.exe -servers=10 -requests=7
LSPQ_RA.exe -servers=10 -requests=8
LSPQ_RA.exe -servers=10 -requests=9
LSPQ_RA.exe -servers=10 -requests=10
cd ../Q1.3;
LSPQ_RA.exe -servers=10 -requests=1
LSPQ_RA.exe -servers=10 -requests=2
LSPQ_RA.exe -servers=10 -requests=3
LSPQ_RA.exe -servers=10 -requests=4
LSPQ_RA.exe -servers=10 -requests=5
LSPQ_RA.exe -servers=10 -requests=6
LSPQ_RA.exe -servers=10 -requests=7
LSPQ_RA.exe -servers=10 -requests=8
LSPQ_RA.exe -servers=10 -requests=9
LSPQ_RA.exe -servers=10 -requests=10