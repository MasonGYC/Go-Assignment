go build fautTolerantIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
# 2. a. primary CM fails at a random time
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=1 -rejoin=0 -fail_backup=0 
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=1 -rejoin=0 -fail_backup=0
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=1 -rejoin=0 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=1 -rejoin=0 -fail_backup=0 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=1 -rejoin=0 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=1 -rejoin=0 -fail_backup=0 

# 2. b. primary CM restarts after the failure
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=1 -rejoin=1 -fail_backup=0 
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=1 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=1 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=1 -rejoin=1 -fail_backup=0 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=1 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=1 -rejoin=1 -fail_backup=0 

# 3.1 primary fails and restarts 2 times
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=2 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=2 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=2 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=2 -rejoin=1 -fail_backup=0 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=2 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=2 -rejoin=1 -fail_backup=0 

# 3.2 primary fails and restarts 3 times
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=3 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=3 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=3 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=3 -rejoin=1 -fail_backup=0 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=3 -rejoin=1 -fail_backup=0
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=3 -rejoin=1 -fail_backup=0 

# 4.1 both fails and restarts 2 times
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=2 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=2 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=2 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=2 -rejoin=1 -fail_backup=1 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=2 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=2 -rejoin=1 -fail_backup=1 

# 4.2 both fails and restarts 3 times
./fautTolerantIvy.exe -servers=10 -request_page=3 -faults=3 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=3 -faults=3 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=10 -request_page=5 -faults=3 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=5 -faults=3 -rejoin=1 -fail_backup=1 
./fautTolerantIvy.exe -servers=10 -request_page=10 -faults=3 -rejoin=1 -fail_backup=1
./fautTolerantIvy.exe -servers=15 -request_page=10 -faults=3 -rejoin=1 -fail_backup=1 