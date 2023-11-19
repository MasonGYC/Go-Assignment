# Go Assignment 2
50.041 Distributed Systems and Computing Go Assignment 2  
Name: Guo Yuchen  
Student ID: 1004885  

# Question 1
`LSPQ`: Lamport Shared Priority Queue  
`LSPQ_RA`: Lamport Shared Priority Queue with Ricart and Agrawala Optimization  
`VOTING`: Voting protocol with deadlock prevention 

## Compilation
To build:
```
go build .\LSPQ.go .\logger.go .\message.go .\PriorityQueue.go .\server.go
go build .\LSPQ_RA.go .\logger.go .\message.go .\PriorityQueue.go .\server.go
go build .\voting.go .\logger.go .\message.go .\PriorityQueue.go .\server.go 
```

To execute：
```
.\LSPQ.exe -servers=? -requests=?
.\LSPQ_RA.exe -servers=? -requests=?
.\voting.exe -servers=? -requests=?
```  
- `-servers`: `int`, indicates the number of clients.  
- `-requests`: `int`, indicates the number of concurrent requests to make.   

## External package
`log` : used for output logging and debugging purpose.   
`container/heap`: used to construct priority queue.  

## Implementation
1. Priority has 2 metrics: 
	1. lamport scalar clock (smaller clock -> higher priority)
	2. node id (higher id -> higher priority)
2. Each server can only make one request at one time. If it doesn't receive all replies, it won't start making new request. 
3. Clocks are updated before the execution of a certain action.
4. If no message is received after 5 seconds, the server will stop. (timeout = 5 seconds)
5. lamport clock is used to avoid deadlock in comparison the timestamp.

## Output interpretation
The `logs.txt` in the folder contains the sample outputs with 1 - 10 clients, while all of them make request concurrently. Refer to the files for more logs.

### Sample outputs with 2 clients
**LSPQ (Lamport Shared Priority Queue)**
```log
10:33:13 server.go:97: Server 1 make request at clock 1.
10:33:13 server.go:97: Server 0 make request at clock 1.
10:33:13 server.go:107: Server 1 requests to 0 at clock 3.
10:33:13 server.go:49: Server 0 received Server 1's request at clock 1.
10:33:13 server.go:107: Server 0 requests to 1 at clock 2.
10:33:13 server.go:180: Server 0 pushed req from 1 to queue.
10:33:13 server.go:49: Server 1 received Server 0's request at clock 2.
10:33:13 server.go:180: Server 1 pushed req from 0 to queue.
10:33:13 server.go:55: Server 1 received Server 0's reply at clock 3.
10:33:13 server.go:208: Server 1's reply_counter is 0.
10:33:13 server.go:127: Server 0 replys to 1 at clock 3.
10:33:13 server.go:188: Server 1 holds reply to 0.
10:33:13 server.go:247: Server 1 clear holding reply from 0 at 1.
10:33:13 server.go:127: Server 1 replys to 0 at clock 5.
10:33:13 server.go:55: Server 0 received Server 1's reply at clock 5.
10:33:13 server.go:208: Server 0's reply_counter is 0.
10:33:15 server.go:163: 2023-11-19 10:33:15.3037356 +0800 CST m=+2.027965201
10:33:15 server.go:168: Server 1 has finished cs execution.
10:33:15 server.go:147: Server 1 release to 0 at clock 7.
10:33:15 server.go:61: Server 0 received Server 1's release at clock 7.
10:33:15 server.go:222: Server 0 has poped req from queue.
10:33:15 server.go:224: Server 0 waiting_for_reply_at_clock is: 1
10:33:15 server.go:225: Server 0 reply_counter is: 0
10:33:17 server.go:163: 2023-11-19 10:33:17.3216252 +0800 CST m=+4.045854801
10:33:17 server.go:168: Server 0 has finished cs execution.
10:33:17 server.go:147: Server 0 release to 1 at clock 10.
10:33:17 server.go:61: Server 1 received Server 0's release at clock 10.
10:33:17 server.go:222: Server 1 has poped req from queue.
10:33:17 server.go:224: Server 1 waiting_for_reply_at_clock is: -1
10:33:17 server.go:225: Server 1 reply_counter is: 1
10:33:22 LSPQ.go:61: Elapsed time:  4.0298504s
```

**LSPQ_RA (Lamport Shared Priority Queue with Ricart and Agrawala Optimization)**
```log
09:44:41 server.go:97: Server 1 make request at clock 1.
09:44:41 server.go:107: Server 1 requests to 0 at 1.
09:44:41 server.go:55: Server 0 received Server 1's request at clock 1.
09:44:41 server.go:180: Server 0 has req from server 0 at clock 1 at head of queue.
09:44:41 server.go:127: Server 0 replys to 1 at 3.
09:44:41 server.go:97: Server 0 make request at clock 1.
09:44:41 server.go:61: Server 1 received Server 0's reply at clock 3.
09:44:41 server.go:107: Server 0 requests to 1 at 3.
09:44:41 server.go:55: Server 1 received Server 0's request at clock 3.
09:44:41 server.go:180: Server 1 has req from server 1 at clock 1 at head of queue.
09:44:41 server.go:191: Server 1 has pushed req from 0 to queue .
09:44:43 server.go:165: 2023-11-19 09:44:43.4381443 +0800 CST m=+2.024096501
09:44:43 server.go:170: Server 1 has finished cs execution.
09:44:43 server.go:146: Server 1 Poped reuqest from 1.
09:44:43 server.go:127: Server 1 replys to 0 at 8.
09:44:43 server.go:146: Server 1 Poped reuqest from 0.
09:44:43 server.go:61: Server 0 received Server 1's reply at clock 8.
09:44:45 server.go:165: 2023-11-19 09:44:45.4523239 +0800 CST m=+4.038276101
09:44:45 server.go:170: Server 0 has finished cs execution.
09:44:45 server.go:146: Server 0 Poped reuqest from 0.
09:44:48 LSPQ_RA.go:60: Elapsed time:  4.0206916s
```
**Voting Protocol**
```log
19:02:08 server.go:110: Server 1 made request at 1.
19:02:08 server.go:140: Server 1 requests to 0 at 1.
19:02:08 server.go:110: Server 0 made request at 1.
19:02:08 server.go:140: Server 0 requests to 0 at 1.
19:02:08 server.go:63: Server 0 received Server 1's request at clock 1.
19:02:08 server.go:63: Server 0 received Server 0's request at clock 1.
19:02:08 server.go:140: Server 1 requests to 1 at 1.
19:02:08 server.go:63: Server 1 received Server 1's request at clock 1.
19:02:08 server.go:163: Server 1 vote to 1 at 3.
19:02:08 server.go:69: Server 1 received Server 1's vote at clock 3.
19:02:08 server.go:140: Server 0 requests to 1 at 1.
19:02:08 server.go:163: Server 0 vote to 1 at 4.
19:02:08 server.go:63: Server 1 received Server 0's request at clock 4.
19:02:08 server.go:69: Server 1 received Server 0's vote at clock 3.
19:02:08 server.go:284: Server 1 received_votes [1 0].
19:02:08 server.go:224: Server 1 started cs execution.
19:02:10 server.go:235: 2023-11-19 19:02:10.2421504 +0800 CST m=+2.021503101
19:02:10 server.go:244: Server 1 has finished cs execution.
19:02:10 server.go:210: Server 1 release to 1 at 8.
19:02:10 server.go:210: Server 1 release to 0 at 10.
19:02:10 server.go:75: Server 0 received Server 1's release at clock 8.
19:02:10 server.go:163: Server 1 vote to 0 at 10.
19:02:10 server.go:75: Server 1 received Server 1's release at clock 8.
19:02:10 server.go:163: Server 0 vote to 0 at 11.
19:02:10 server.go:69: Server 0 received Server 1's vote at clock 10.
19:02:10 server.go:284: Server 0 received_votes [1 0].
19:02:10 server.go:69: Server 0 received Server 0's vote at clock 10.
19:02:10 server.go:224: Server 0 started cs execution.
19:02:12 server.go:235: 2023-11-19 19:02:12.264371 +0800 CST m=+4.043723701
19:02:12 server.go:244: Server 0 has finished cs execution.
19:02:12 server.go:210: Server 0 release to 1 at 14.
19:02:12 server.go:210: Server 0 release to 0 at 14.
19:02:12 server.go:75: Server 1 received Server 0's release at clock 14.
19:02:12 server.go:75: Server 0 received Server 0's release at clock 14.
19:02:17 voting.go:59: Elapsed time:  4.0278337s
```

# Question 2
The performance is measured on 1-10 servers making concurrent requests. The time is calculated from teh first request to all servers finished executing critical sections.
- The `performance_logs.txt` in the respective folder contains the time measured for each protocol.  

## Performance table:
unit: second  

| Number of servers | LSPQ | LSPQ_RA | VOTING |
| ---------- | ---------- | ---------- | ---------- |
| 1 | 2.0231424 |  2.0067398 |  2.0024699 |
| 2 | 4.0298504 |  4.0206916 |  4.0278337 |
| 3 | 6.0321942 |  6.0275846 |  6.0414826 |
| 4 | 8.0714022 |  8.0280027 |  8.0508024 | 
| 5 | 10.1354729 |  10.0886073 |  10.0767158 |
| 6 | 12.0911982 |  12.1181105 |  12.0960925 |
| 7 | 14.218373 |  14.1397041 |  14.1663692 |
| 8 | 16.25372 |  16.1298288 |  16.13023 |
| 9 | 18.1930937 |  18.1006392 |  18.1623731 |
| 10 | 20.2354588 |  20.1437053 |  20.1998621 |

From the table, we can observe that generally, LSPQ is slower than VOTING, and VOTING is slower than LSPQ_RA.  