# Question 2
## Compilation

To build: `go build .\run.go .\common.go .\logger.go .\message.go .\node.go`  

To execute：`.\run.exe -nodes=3 -sync=1 -timeout=6`.
- -nodes `int` number of nodes (default 3)  
- -sync `int` the time interval to sync in second (default 1)  
- -timeout `int` 2T(m) + T(p) in second (default 6)  


## External package
`log` : used for output logging and debugging purpose.

## Implementation

### Message
Message is defined to be a struct with 4 fields: 
 
- `sender_id`: `int`, the id of the sender machine
- `message_type`: `string`, the type of the message, including:  
 	- `SYNC` Coordinator sync data with replica.  
    - `ELECT` A node elect itself to be the new coordinator.  
	- `ACK` A node refuse one's self-elect request.   
	- `VICTORY` A node broadcast that it is the now coordinator. 
	- `STOP` To stop a goroutine.     
- `message`: `any`, usually a sentence about the content, which is readable by human.
- `timestamp`: `time.Time`, the physical clock of the message. currently not used in the program.

### Node
Each node has a state, a role, 4 channels. 
- Role:
    - `REPLICA` 
    - `COORDINATOR`
- State:
    - `NORMAL` node performs normal work, not during election
	- `SELF_ELECTING` node is sending self-electing messages
	- `BROADCATING` new coordinator broadcasting victory
	- `DOWN` node fails
- Channels:
	- `ch_sync` for SYNC message  
	- `ch_elect` for ELECT, ACK, VICTORY message  
	- `ch_stop_elect` for STOP mesage to stop self_elect() or broadcast_victory()  
	- `ch_role_switch` for STOP mesage to switch between `COORDINATOR` and `REPLICA`   

### Synchonization and Election
- **Synchonization**  
 If Coordinator alive, send SYNC message at certain interval to Replica. If Replica receives SYNC message, update its data; if not for a timeout, start a election to elect itself as the Coordinator.

- **Election by bully algorithm**  
 A node sends ELECT message to nodes with higher ids. If receives ACK or STOP, go back to NORMAL state and stops election. If not for a certain time, broadcasts its victory to all. After broadcasting, switch role to Coordinator and start sending SYNC messages.

### Failure Implementation
 Set node state to `DOWN` and stop preocessing received message,sending message, and electing. It can still receive messages but won't take care of it.

 For implementation details, refer to the code.

## Output
The `logs.txt` in each folder contains the sample outputs with 3 nodes, and the coordinator fails after 4 seconds.  

**Process:**  
Normal data sync between node 1,2,3->   
Coordinator fails ->   
Node 0 and 1 both detect failure and start election ->     
Node 1 refuse Node 0's request ->  
Node 0 ends election ->  
Node 1 won election and become coordinator ->  
Normal data sync between node 1,2  

### Quetsion checklist
- [DONE] Multiple GO routines start the election process simultaneously. Thus the worst and best case are covered as well.
	- set `failcoor` = true when execute
	- sample output: `logs_coor_fails.txt` contains 3 and 4 nodes cases

- [DONE] Coordinator silently leaves the network. 
	- set `failcoor` = true when execute
	- sample output: `logs_coor_fails.txt` contains 3 and 4 nodes cases

- [DONE] Replica silently leaves the network.
	- set `failrep` = true when execute
	- sample output: `logs_rep_fails.txt` contains 3 and 4 nodes cases

- [x] The newly elected coordinator fails while announcing.
	- set `flcrel` = true and `fail_coordinator`= true when execute


- [x] The failed node is not the newly elected coordinator.
	- set `flrpel` = true and `fail_coordinator`= true when execute



## Others
### Issues
1. The current program only supports <=3 nodes running, sometimes up to 4, due to limited computing power of my laptop. Channels get stuck when too many messages arrive. Thus the next step would be to create separate channels for ELECT, ACK, VICTORY messages.
2. Some of the unexpected behaviours are not considered due to the limited computing power thus not able to debug. especially when there are too many nodes, 

### Assumptions
1. Network is reliable.
2. Network is asynchoronous.
3. Channels won't congest. (But it actually happens when node number increases.)
