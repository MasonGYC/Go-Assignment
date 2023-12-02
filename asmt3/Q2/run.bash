
go build fautTolerantIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
./fautTolerantIvy.exe -servers=1 -requests=1 
./fautTolerantIvy.exe -servers=2 -requests=2
./fautTolerantIvy.exe -servers=3 -requests=3
./fautTolerantIvy.exe -servers=4 -requests=4
./fautTolerantIvy.exe -servers=5 -requests=5
./fautTolerantIvy.exe -servers=6 -requests=6
./fautTolerantIvy.exe -servers=7 -requests=7
./fautTolerantIvy.exe -servers=8 -requests=8
./fautTolerantIvy.exe -servers=9 -requests=9
./fautTolerantIvy.exe -servers=10 -requests=10