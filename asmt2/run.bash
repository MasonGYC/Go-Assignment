cd Q1.1;
go build LSPQ.go logger.go message.go PriorityQueue.go server.go
LSPQ.exe -servers=1 -requests=1
LSPQ.exe -servers=2 -requests=2
LSPQ.exe -servers=3 -requests=3
LSPQ.exe -servers=4 -requests=4
LSPQ.exe -servers=5 -requests=5
LSPQ.exe -servers=6 -requests=6
LSPQ.exe -servers=7 -requests=7
LSPQ.exe -servers=8 -requests=8
LSPQ.exe -servers=9 -requests=9
LSPQ.exe -servers=10 -requests=10
