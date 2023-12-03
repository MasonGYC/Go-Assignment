# Experiments 
## 1. Without any faults
**Requirement**: Compare the performance of the basic version of Ivy protocol and the new fault tolerant version using requests from at least 10 clients.

The simulation code is as follows. There will be `num_servers` servers making `num_request_page` requests each, therefore in total `num_servers * num_request_page` requests in one round. The request is raised every 1 second. 

```go
for i := 0; i < *num_servers; i++ {
    for j := 0; j < *num_request_page; j++ {
        wg.Add(1)
        go func(i int) {
            if j%2 == 1 {
                servers[i].read(j)
            } else {
                servers[i].write(j)
            }
            wg.Done()
        }(i)
        time.Sleep(time.Second)
    }
}
```
The measured time(in second) is as follows:
| Number of servers* | Number of requests* | Vanila Ivy | Fault Tolerant Ivy | 
|----------|----------|----------|----------|
| 10 | 3 | 14.73 | 14.74 |
| 15 | 3 | 22.46 | 22.40 |
| 10 | 5 | 24.95 | 24.92 |
| 15 | 5 | 37.10 | 40.58 |
| 10 | 10 | 50.02 | 49.96 |
| 15 | 10 | 75.88 | 75.88 |

*Number of servers = num_servers
*Number of requests = num_request_page

From the data we can see that the time is roughly the same, since the fault tolerance version shares almost the same process as the basic version, other than the data sync and heatbeat, which does not cause heavy overhead.

## 2. One CM fails only once
**Requirement**: Simulate 2 cases:  
- a) when the primary CM fails at a random time,   
- b) when the primary CM restarts after the failure.   

The simulation code is as follows:
```go
go func() {
    // simulate primary down
    primaryManager.down()

    if *rejoin_primary {
        // simulate primary rejoin
        time.Sleep(2 * time.Second)
        primaryManager.rejoin()
    }
}()
```

The measured time(in second) is as follows: (the last 2 columns are from this experiment)  
| Number of servers | Number of requests |  Primary fails | Primary fails and restarts | 
|----------|----------|----------|----------|
| 10 | 3 | 14.71 | 14.74 |
| 15 | 3 | 22.31 | 22.33 |
| 10 | 5 | 24.81 | 24.83 |
| 15 | 5 | 39.02 | 8.07  |
| 10 | 10 | 49.65 | 49.74 |
| 15 | 10 | 75.39 | 46.10 |

From the data we can see that the time is roughly the same, since when primary restarts, it only takes a very short time to get the control back. And compared with the case without any faults, the time generally increased since backup server needs to broadcast to all servers about its primary role, and some messages may get lost during the transfering period. And compared with fault-free version, the time increased due to the handling-over by the backup manager.

## 3. Multiple faults for primary CM
**Requirement**: Primary CM fails and restarts multiple times.

The simulation code is as follows:
```go
go func() {
    for i := 0; i < *num_faults_primary; i++ {
        // simulate primary down
        primaryManager.down()

        // simulate primary rejoin
        time.Sleep(2 * time.Second)
        primaryManager.rejoin()

        if *fail_backup {
            // simulate backup down
            backupManager.down()

            // simulate backup rejoin
            time.Sleep(2 * time.Second)
            backupManager.rejoin()
        }
        // rest a while, don't fail so frequently
        time.Sleep(500 * time.Millisecond)
    }
}()
```
Let primary fails and restarts 2, 3 times.
The measured time(in second) is as follows:
| Number of servers | Number of requests | 2 | 3 |
|----------|----------|----------| ----------|
| 10 | 3 | 14.69 | 14.65 |
| 15 | 3 | 22.31 | 22.18 |
| 10 | 5 | 24.80 | 24.68 |
| 15 | 5 | 37.49 | 37.27 |
| 10 | 10 | 49.61 | 49.32 |
| 15 | 10 | 75.32 | 74.89 |

From the data we can see that the time is roughly the same, since the control transfer is fast enough, and lost message won't be retransmissioned by design. Thus it won't affect the general performance in terms of time, but the number of requests missing will increase.

## 4. Multiple faults for primary CM and backup CM 
**Requirement**: Both primary CM and backup CM fail and restart multiple times. 

The simulation code is the sam as (3), but set `fail_backup=true`.  
Let both fail and restart 2, 3 times.  
The measured time(in second) is as follows:  
| Number of servers | Number of requests | 2 | 3 |
|----------|----------|----------| ----------|
fail 2 times

fail 3 times
| 10 | 3 | 19075749900 |
| 15 | 3 | 20085428600 |
| 10 | 5 | 19070842400 |
| 15 | 5 | 19093296200 |
| 10 | 10 | 19071587300 |
| 15 | 10 | 19087787600 |