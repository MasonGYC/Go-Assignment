# go build voting.go logger.go message.go queue.go server.go clock.go
# voting.exe -servers=1 -requests=1
# voting.exe -servers=2 -requests=2
# voting.exe -servers=3 -requests=3
# voting.exe -servers=4 -requests=4
# voting.exe -servers=5 -requests=5
# voting.exe -servers=6 -requests=6
# voting.exe -servers=7 -requests=7
# voting.exe -servers=8 -requests=8
# voting.exe -servers=9 -requests=9
# voting.exe -servers=10 -requests=10

go build LSPQ_RA.go logger.go message.go PriorityQueue.go server.go
LSPQ_RA.exe -servers=1 -requests=1
LSPQ_RA.exe -servers=2 -requests=2
LSPQ_RA.exe -servers=3 -requests=3
LSPQ_RA.exe -servers=4 -requests=4
LSPQ_RA.exe -servers=5 -requests=5
LSPQ_RA.exe -servers=6 -requests=6
LSPQ_RA.exe -servers=7 -requests=7
LSPQ_RA.exe -servers=8 -requests=8
LSPQ_RA.exe -servers=9 -requests=9
LSPQ_RA.exe -servers=10 -requests=10