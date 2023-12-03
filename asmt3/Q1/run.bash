go build vanilaIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
./vanilaIvy.exe -servers=10 -requests=3
./vanilaIvy.exe -servers=15 -requests=3 
./vanilaIvy.exe -servers=10 -requests=5
./vanilaIvy.exe -servers=15 -requests=5 
./vanilaIvy.exe -servers=10 -requests=10
./vanilaIvy.exe -servers=15 -requests=10 