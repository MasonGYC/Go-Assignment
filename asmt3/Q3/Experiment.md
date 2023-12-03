# Experiments 
## 1. Without any faults
**Requirement**: Compare the performance of the basic version of Ivy protocol and the new fault tolerant version using requests from at least 10 clients.

The simulation code is as follows. There will be `num_servers` servers making 3 requests each, therefore in total `num_servers * num_request_page` requests in one round. The request is raised every 1 second. 

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
| Number of servers | Number of requests | Vanila Ivy | Fault Tolerant Ivy | 
|----------|----------|----------|----------|
| 10 | 3 | 14727643600 |
| 15 | 3 | 22462469800 |
| 20 | 3 | 30153258300 |
| 10 | 5 | 24947046300 |
| 15 | 5 | 37097288300 |
| 20 | 5 | 49893654800 |
| 10 | 10 | 50022001100 |
| 15 | 10 | 75880231800 |
| 20 | 10 | 101394502100 |
14.7276436
22.4624698
30.1532583
24.9470463
37.0972883
49.8936548
50.0220011
75.8802318
101.3945021
ft:
| 10 | 3 | 14735053800 |
| 15 | 3 | 22396158500 |
| 20 | 3 | 31091008800 |
| 10 | 5 | 24922633000 |
| 15 | 5 | 40583409400 |
| 20 | 5 | 60777452000 |
| 10 | 10 | 49956169500 |
| 15 | 10 | 75875247100 |
| 20 | 10 | 101086242400 |
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
| Number of servers | Vanila Ivy | Fault Tolerant Ivy | Primary fails | Primary fails and restarts | 
|----------|----------|----------|----------|----------|
| 11 | 61.78 | 61.78 | 61.89 | 61.90 |
| 12 | 72.79 | 72.75 | 73.06 | 72.88 |
| 13 | 85.93 | 85.80 | 86.26 | 85.98 |

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
| Number of servers | 1 | 2 | 3 |
|----------|----------|----------| ----------|
| 11 |  61.90 | 62.12 |
| 12 |  72.88 | 73.50 |
| 13 |  85.98 | 86.54 |

- `NA` is caused due to message lost. A 

## 4. Multiple faults for primary CM and backup CM 
**Requirement**: Both primary CM and backup CM fail and restart multiple times. 

The simulation code is the sam as (3), but set `fail_backup=true`.  
Let both fail and restart 2, 3 times.  
The measured time(in second) is as follows:  
| Number of servers | 1 | 2 | 3 |
|----------|----------|----------| ----------|
| 11 |  61.90 |  |
| 12 |  72.88 |  |
| 13 |  85.98 |  |