# Question 1
## Compilation
For quetsion 1 and 2:  
To build: `go build .\lamport.go`
To execute：`.\lamport.exe -clients=15`. `clients` indicates the number of clients. The default number is 15.

For quetsion 3, run  
To build: `go build .\vector.go`
To execute：`.\vector.exe -clients=15`. `clients` indicates the number of clients. The default number is 15.

## External package
`log` : used for output logging and debugging purpose.

## Implementation

1. In the program, Message is defined to be a struct with 3 fields: 
- `sender_id`: `int`, the id of the sender client/server
- `message`: `int`, the i-th message being sent by a client
- `clock`: `int` for lamport clock, `[]int` for vector clock
2. Clients send a message every 5 seconds.
3. There's a 50-50 chance of dropping or forwarding received message for the server. (flag can be 0 or 1, if 0, forward message)
4. Clock is updated every time an action (i.e. receiving message, send message) is performed.
5. In `vector.go`, causality violation is checked every time a message is received (`clockIsGreaterThan(c1, c2)`) If the local vector clock of the receiving machine is more than the vector clock of the message, then a potential causality violation is detected.  

## Output interpretation
The `logs.txt` in each folder contains the sample outputs with 15 clients. Refer to the files for more logs.

### Lamport
```log
08:59:22 ===============START===============
08:59:22 Client 7's clock:  1.
08:59:22 Client 7 sending the 0-th message at clock 1.
08:59:22 Server received Client 7's 0-th message of clock 1.
08:59:22 Server's clock: 2.
08:59:22 Server forward to Client 0 Client 7's 0-th message.
...
08:59:22 Client 3 sending the 0-th message at clock 1.
08:59:22 Client 0 received Server's message. Sent from Client 7 at clock 2.
...
```
### Vector
```log
09:03:31 ===============START===============
09:03:31 Client 7's clock: [0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0].
...
09:03:31 Client 3 sending the 0-th message at clock [0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0].
09:03:31 Client 9's clock: [0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0].
09:03:31 Client 9 sending the 0-th message at clock [0 0 0 0 0 0 0 0 0 1 0 0 0 0 0 0].
...
09:03:31 Client 6 sending the 0-th message at clock [0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0].
09:03:31 Client 14's clock: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 0].
09:03:31 Client 7 sending the 0-th message at clock [0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0].
...
09:03:31 Server received Client 3's 0-th message of clock [0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0].
09:03:31 Server's clock: [0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 1].
09:03:31 Drop Client 3's 0-th message of clock [0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0].

...
09:03:31 Server's clock: [0 0 0 1 0 0 1 1 0 1 1 1 0 0 0 6].
09:03:31 Server forward to Client 0 Client 10's 0-th message.
09:03:31 Server forward to Client 1 Client 10's 0-th message.
...
09:03:31 Server forward to Client 14 Client 0's 1-th message.
09:03:31 Client 13 received Server's message. Sent from Client 0 at clock [10 8 5 1 6 5 1 1 2 1 1 1 4 3 6 16].
...
09:03:31 Client 14 received Server's message. Sent from Client 0 at clock [10 8 5 1 6 5 1 1 2 1 1 1 4 3 6 16].
09:03:31 Potential violation detected!
09:03:31 Local clock: [10 8 5 1 6 5 1 1 2 1 1 1 4 8 6 16].
09:03:31 Message clock: [10 8 5 1 6 5 1 1 2 1 1 1 4 3 6 16].
...
```

## Others
### Assumptions
1. Network is reliable.
2. Network is asynchoronous.
3. Channels won't congest. (But it actually happens when node number increases.)
