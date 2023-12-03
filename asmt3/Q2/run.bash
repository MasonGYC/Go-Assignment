
go build fautTolerantIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=1 -rejoin=1 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=1 -rejoin=0 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=20 -request_page=3 -faults=1 -rejoin=0 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=1 -rejoin=0 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=1 -rejoin=0 -fail_backup=0  
# ./fautTolerantIvy.exe -servers=20 -request_page=5 -faults=1 -rejoin=0 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=1 -rejoin=0 -fail_backup=0 
# ./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=1 -rejoin=0 -fail_backup=0  
# ./fautTolerantIvy.exe -servers=20 -request_page=10 -faults=1 -rejoin=0 -fail_backup=0 