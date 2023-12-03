# Go Assignment 3
50.041 Distributed Systems and Computing Go Assignment 3  
Name: Guo Yuchen  
Student ID: 1004885  

# Question 1
`vanilaIvy`: Basic verison of Ivy   
`fautTolerantIvy`: Fault tolerance verison of Ivy  

## Compilation
To build:
```
go build vanilaIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
go build fautTolerantIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go
```

To execute：
```
./vanilaIvy.exe -servers=? -requests=? 

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

### Sample outputs
10 severs, 8 requests, primary down once and rejoin once
```log
17:14:07 manager.go:237: Manager -2 received heartbeat from -1 at clock 0.
17:14:07 manager.go:134: Manager -1 sent heartbeat to -2 at clock 2.
17:14:07 server.go:32: Server 7 started listening...
17:14:07 manager.go:237: Manager -1 received heartbeat from -2 at clock 2.
17:14:08 manager.go:134: Manager -1 sent heartbeat to -2 at clock 4.
17:14:08 manager.go:237: Manager -2 received heartbeat from -1 at clock 1.
17:14:08 manager.go:237: Manager -1 received heartbeat from -2 at clock 4.
17:14:08 server.go:93: Server 0 wants to read page 3...
17:14:08 server.go:104: Server 0 request to read page 3 to manager -1 at clock 1.
17:14:08 server.go:77: Server 0 resends request to manager -1 at clock 1.
17:14:08 manager.go:53: Manager -1 received Server 0's message at clock 1.
17:14:08 manager.go:61: Manager -1 received Server 0's read request at clock 1.
17:14:08 manager.go:205: Manager -1 forward READ request for page 3 at clock 6.
17:14:08 server.go:43: Server 3 received Manager -1's read forward at clock 6.
17:14:08 server.go:58: Server 0 received Server 3's sent page at clock 8.
17:14:08 server.go:141: Server 3 sent page 3 to Server 0 at clock 8.
17:14:08 server.go:180: Server 0 is reading page 3...
17:14:08 server.go:155: Server 0 sent READ confirm for page 3 to manager -1 at clock 10.
17:14:08 manager.go:53: Manager -1 received Server 0's message at clock 10.
17:14:08 manager.go:70: Manager -1 received Server 0's read confirm at clock 10.
17:14:08 manager.go:339: Manager -1 adds server 0 to page 3's copyset.
17:14:08 manager.go:123: Manager -2 updated records.
17:14:08 manager.go:364: Manager -1 is down.
17:14:08 server.go:93: Server 1 wants to read page 6...
17:14:08 server.go:104: Server 1 request to read page 6 to manager -1 at clock 1.
17:14:08 server.go:77: Server 1 resends request to manager -1 at clock 1.
17:14:08 manager.go:53: Manager -1 received Server 1's message at clock 1.
17:14:09 server.go:93: Server 2 wants to read page 6...
17:14:09 server.go:104: Server 2 request to read page 6 to manager -1 at clock 1.
17:14:09 server.go:77: Server 2 resends request to manager -1 at clock 1.
17:14:09 manager.go:53: Manager -1 received Server 2's message at clock 1.
17:14:09 server.go:93: Server 3 wants to read page 0...
17:14:09 server.go:104: Server 3 request to read page 0 to manager -1 at clock 9.
17:14:09 server.go:77: Server 3 resends request to manager -1 at clock 9.
17:14:09 manager.go:53: Manager -1 received Server 3's message at clock 9.
17:14:10 manager.go:100: Manager -1 detects the other manager is down.
17:14:10 manager.go:100: Manager -2 detects the other manager is down.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 0 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 1 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 2 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 3 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 4 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 5 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 6 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 7 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 8 at clock 2.
17:14:10 manager.go:107: Declare manager -2 to be pri to server 9 at clock 2.
17:14:10 server.go:67: Server 9 received Server -2's sent page at clock 2.
17:14:10 server.go:67: Server 3 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 0's new primary manager updated: -2.
17:14:10 server.go:210: Server 3's new primary manager updated: -2.
17:14:10 server.go:210: Server 8's new primary manager updated: -2.
17:14:10 server.go:67: Server 0 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 2's new primary manager updated: -2.
17:14:10 server.go:67: Server 4 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 4's new primary manager updated: -2.
17:14:10 server.go:67: Server 5 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 5's new primary manager updated: -2.
17:14:10 server.go:67: Server 6 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 6's new primary manager updated: -2.
17:14:10 server.go:67: Server 7 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 7's new primary manager updated: -2.
17:14:10 server.go:67: Server 8 received Server -2's sent page at clock 2.
17:14:10 server.go:67: Server 1 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 1's new primary manager updated: -2.
17:14:10 server.go:67: Server 2 received Server -2's sent page at clock 2.
17:14:10 server.go:210: Server 9's new primary manager updated: -2.
17:14:10 server.go:93: Server 4 wants to read page 8...
17:14:10 server.go:104: Server 4 request to read page 8 to manager -2 at clock 4.
17:14:10 server.go:77: Server 4 resends request to manager -2 at clock 4.
17:14:10 manager.go:53: Manager -2 received Server 4's message at clock 4.
17:14:10 manager.go:61: Manager -2 received Server 4's read request at clock 4.
17:14:10 manager.go:205: Manager -2 forward READ request for page 8 at clock 6.
17:14:10 server.go:43: Server 8 received Manager -2's read forward at clock 6.
17:14:10 server.go:141: Server 8 sent page 8 to Server 4 at clock 8.
17:14:10 server.go:180: Server 4 is reading page 8...
17:14:10 server.go:155: Server 4 sent READ confirm for page 8 to manager -2 at clock 10.
17:14:10 manager.go:53: Manager -2 received Server 4's message at clock 10.
17:14:10 manager.go:70: Manager -2 received Server 4's read confirm at clock 10.
17:14:10 manager.go:339: Manager -2 adds server 4 to page 8's copyset.
17:14:10 server.go:58: Server 4 received Server 8's sent page at clock 8.
17:14:10 manager.go:368: Manager -1 rejoined.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 0 at clock 0.
17:14:10 manager.go:134: Manager -1 sent heartbeat to -2 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 1 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 2 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 3 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 4 at clock 2.
17:14:10 server.go:67: Server 2 received Server -1's sent page at clock 2.
17:14:10 server.go:67: Server 1 received Server -1's sent page at clock 2.
17:14:10 server.go:67: Server 0 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 0's new primary manager updated: -1.
17:14:10 server.go:210: Server 3's new primary manager updated: -1.
17:14:10 manager.go:237: Manager -2 received heartbeat from -1 at clock 12.
17:14:10 manager.go:123: Manager -1 updated records.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 5 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 6 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 7 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 8 at clock 2.
17:14:10 manager.go:374: Declare manager -1 to be pri to server 9 at clock 2.
17:14:10 server.go:67: Server 9 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 9's new primary manager updated: -1.
17:14:10 server.go:67: Server 4 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 1's new primary manager updated: -1.
17:14:10 server.go:210: Server 2's new primary manager updated: -1.
17:14:10 server.go:210: Server 4's new primary manager updated: -1.
17:14:10 server.go:67: Server 3 received Server -1's sent page at clock 2.
17:14:10 server.go:67: Server 5 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 5's new primary manager updated: -1.
17:14:10 server.go:67: Server 8 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 8's new primary manager updated: -1.
17:14:10 server.go:67: Server 6 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 6's new primary manager updated: -1.
17:14:10 server.go:67: Server 7 received Server -1's sent page at clock 2.
17:14:10 server.go:210: Server 7's new primary manager updated: -1.
17:14:10 server.go:114: Server 5 wants to write page 5...
17:14:11 server.go:93: Server 6 wants to read page 9...
17:14:11 server.go:104: Server 6 request to read page 9 to manager -1 at clock 5.
17:14:11 server.go:77: Server 6 resends request to manager -1 at clock 5.
17:14:11 manager.go:53: Manager -1 received Server 6's message at clock 5.
17:14:11 manager.go:61: Manager -1 received Server 6's read request at clock 5.
17:14:11 manager.go:205: Manager -1 forward READ request for page 9 at clock 7.
17:14:11 server.go:43: Server 9 received Manager -1's read forward at clock 7.
17:14:11 server.go:58: Server 6 received Server 9's sent page at clock 9.
17:14:11 server.go:141: Server 9 sent page 9 to Server 6 at clock 9.
17:14:11 server.go:180: Server 6 is reading page 9...
17:14:11 manager.go:53: Manager -1 received Server 6's message at clock 11.
17:14:11 server.go:155: Server 6 sent READ confirm for page 9 to manager -1 at clock 11.
17:14:11 manager.go:70: Manager -1 received Server 6's read confirm at clock 11.
17:14:11 manager.go:339: Manager -1 adds server 6 to page 9's copyset.
17:14:11 manager.go:123: Manager -2 updated records.
17:14:11 manager.go:134: Manager -1 sent heartbeat to -2 at clock 15.
17:14:11 manager.go:237: Manager -2 received heartbeat from -1 at clock 13.
17:14:11 manager.go:237: Manager -1 received heartbeat from -2 at clock 15.
17:14:11 manager.go:53: Manager -1 received Server 3's message at clock 9.
17:14:11 manager.go:61: Manager -1 received Server 3's read request at clock 9.
17:14:11 manager.go:205: Manager -1 forward READ request for page 0 at clock 17.
17:14:11 server.go:43: Server 0 received Manager -1's read forward at clock 17.
17:14:11 server.go:141: Server 0 sent page 0 to Server 3 at clock 19.
17:14:11 server.go:58: Server 3 received Server 0's sent page at clock 19.
17:14:11 server.go:180: Server 3 is reading page 0...
17:14:11 server.go:155: Server 3 sent READ confirm for page 0 to manager -1 at clock 21.
17:14:11 manager.go:53: Manager -1 received Server 3's message at clock 21.
17:14:11 manager.go:70: Manager -1 received Server 3's read confirm at clock 21.
17:14:11 manager.go:339: Manager -1 adds server 3 to page 0's copyset.
17:14:11 manager.go:123: Manager -2 updated records.
17:14:11 server.go:120: Server 5 finished writing page 5.
17:14:11 server.go:93: Server 7 wants to read page 9...
17:14:11 server.go:104: Server 7 request to read page 9 to manager -1 at clock 5.
17:14:11 server.go:77: Server 7 resends request to manager -1 at clock 5.
17:14:11 manager.go:53: Manager -1 received Server 7's message at clock 5.
17:14:11 manager.go:61: Manager -1 received Server 7's read request at clock 5.
17:14:11 manager.go:205: Manager -1 forward READ request for page 9 at clock 25.
17:14:11 server.go:43: Server 9 received Manager -1's read forward at clock 25.
17:14:11 server.go:141: Server 9 sent page 9 to Server 7 at clock 27.
17:14:11 server.go:58: Server 7 received Server 9's sent page at clock 27.
17:14:11 server.go:180: Server 7 is reading page 9...
17:14:11 server.go:155: Server 7 sent READ confirm for page 9 to manager -1 at clock 29.
17:14:11 manager.go:53: Manager -1 received Server 7's message at clock 29.
17:14:11 manager.go:70: Manager -1 received Server 7's read confirm at clock 29.
17:14:11 manager.go:339: Manager -1 adds server 7 to page 9's copyset.
17:14:11 manager.go:123: Manager -2 updated records.
17:14:12 manager.go:134: Manager -1 sent heartbeat to -2 at clock 33.
17:14:12 manager.go:237: Manager -2 received heartbeat from -1 at clock 14.
...
17:14:16 fautTolerantIvy.go:127: Elapsed time:  8.5481405s
```

# Q3 
**Discuss whether your fault tolerant version of Ivy still preserves sequential consistency or not.**
Yes. Because on each individual process, operations are still executed in sequence. The central manager, despite susceptible to fault, will serve all write request in a FIFO queue, thus eliminate the possibility of a server reading stale values. And this queue will be synchorinized with the backup manager everytime it gets updated. Thus when the central manager fails, the backup manager takes the same queue and continue execution. There may be message lost during the transferring process, but it does not affect the sequential consistency.