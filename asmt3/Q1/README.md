# implementations:
1. manager holds a priority queue, requests will be stored in queue and executed sequentially. Request with earlier timestamp or higher server id are prioritised.

# records
initially, every server possess one page, the page number of it is its own id

# TODO


go run vanilaIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go


# Logs 
## 2 writes
```log
09:51:01 vanilaIvy.go:24: ===============START===============
09:51:01 server.go:88: Server 1 wants to write page 3...
09:51:01 server.go:99: Server 1 request to write page 3 to manager at clock 0.
09:51:01 server.go:88: Server 0 wants to write page 3...
09:51:01 server.go:99: Server 0 request to write page 3 to manager at clock 0.
09:51:01 manager.go:46: Manager received Server 1's write request at clock 0.
09:51:01 manager.go:46: Manager received Server 0's write request at clock 0.
09:51:01 manager.go:161: Manager pushed server 1's wr req for page 3 at clock 1.
09:51:01 manager.go:161: Manager pushed server 0's wr req for page 3 at clock 2.
09:51:01 manager.go:74: Signal channel received msg  {-1 -1 3 NOTICE_WRITE_REQUEST 0 {0 <nil>} 1}  at clock  3
09:51:01 manager.go:79: record.writing is false at clock 3.
09:51:01 manager.go:74: Signal channel received msg  {-1 -1 3 NOTICE_WRITE_REQUEST 0 {0 <nil>} 0}  at clock  4
09:51:01 manager.go:79: record.writing is true at clock 4.
09:51:01 manager.go:87: Manager poped s1's wr req for page 3.
09:51:01 manager.go:117: Manager forward WRITE request for page 3 at clock 5.
09:51:01 server.go:45: Server 3 received Manager's write forward at clock 5.
09:51:01 server.go:113: Server 3 sent page 3 to Server 1 at clock 7.
09:51:01 server.go:60: Server 1 received Server 3's sent page at clock 7.
09:51:01 server.go:170: Server 1 is writing page 3 at clock 8.
09:51:01 server.go:127: Server 1 sent WRITE confirm for page 3 to manager at clock 9.
09:51:01 manager.go:54: Manager received Server 1's write confirm at clock 9.
09:51:01 manager.go:74: Signal channel received msg  {-1 -1 3 NOTICE_WRITE_REQUEST 9 {0 <nil>} 1}  at clock  11
09:51:01 manager.go:79: record.writing is false at clock 11.
09:51:01 manager.go:87: Manager poped s0's wr req for page 3.
09:51:01 manager.go:117: Manager forward WRITE request for page 3 at clock 12.
09:51:01 server.go:45: Server 1 received Manager's write forward at clock 12.
09:51:01 server.go:113: Server 1 sent page 3 to Server 0 at clock 14.
09:51:01 server.go:60: Server 0 received Server 1's sent page at clock 14.
09:51:01 server.go:170: Server 0 is writing page 3 at clock 15.
09:51:01 server.go:127: Server 0 sent WRITE confirm for page 3 to manager at clock 16.
09:51:01 manager.go:54: Manager received Server 0's write confirm at clock 16.
09:51:11 vanilaIvy.go:82: Elapsed time:  10.0095429s
```
## 2 reads
```log
10:03:34 vanilaIvy.go:24: ===============START===============
10:03:34 server.go:71: Server 1 wants to read page 3...
10:03:34 server.go:82: Server 1 request to read page 3 to manager at clock 0.
10:03:34 server.go:71: Server 0 wants to read page 3...
10:03:34 server.go:82: Server 0 request to read page 3 to manager at clock 0.
10:03:34 manager.go:41: Manager received Server 1's read request at clock 0.
10:03:34 manager.go:41: Manager received Server 0's read request at clock 0.
10:03:34 manager.go:114: Manager forward READ request for page 3 at clock 3.
10:03:34 server.go:40: Server 3 received Manager's read forward at clock 3.
10:03:34 server.go:113: Server 3 sent page 3 to Server 0 at clock 5.
10:03:34 server.go:55: Server 0 received Server 3's sent page at clock 5.
10:03:34 server.go:152: Server 0 is reading page 3...
10:03:34 server.go:127: Server 0 sent READ confirm for page 3 to manager at clock 7.
10:03:34 manager.go:50: Manager received Server 0's read confirm at clock 7.
10:03:34 manager.go:227: Manager adds server 0 to page 3's copyset.
10:03:34 manager.go:114: Manager forward READ request for page 3 at clock 9.
10:03:34 server.go:40: Server 3 received Manager's read forward at clock 9.
10:03:34 server.go:113: Server 3 sent page 3 to Server 1 at clock 11.
10:03:34 server.go:55: Server 1 received Server 3's sent page at clock 11.
10:03:34 server.go:152: Server 1 is reading page 3...
10:03:34 server.go:127: Server 1 sent READ confirm for page 3 to manager at clock 13.
10:03:34 manager.go:50: Manager received Server 1's read confirm at clock 13.
10:03:34 manager.go:227: Manager adds server 1 to page 3's copyset.
10:03:44 vanilaIvy.go:82: Elapsed time:  10.0007147s
```