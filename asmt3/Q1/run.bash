
go build vanilaIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
./vanilaIvy.exe -servers=10 -request_page=3
./vanilaIvy.exe -servers=15 -request_page=3 
./vanilaIvy.exe -servers=20 -request_page=3
./vanilaIvy.exe -servers=10 -request_page=5
./vanilaIvy.exe -servers=15 -request_page=5 
./vanilaIvy.exe -servers=20 -request_page=5
./vanilaIvy.exe -servers=10 -request_page=10
./vanilaIvy.exe -servers=15 -request_page=10 
./vanilaIvy.exe -servers=20 -request_page=10