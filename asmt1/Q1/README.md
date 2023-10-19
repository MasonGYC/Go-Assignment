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

## Output interpretation
PS: output is shown in both terminal and the `logs.txt` file which is at the same directory as the program file.

1. In the program, Message is defined to be a struct with 3 fields: 
- `sender_id`: `int`, the id of the sender client/server
- `message`: `int`, the i-th message being sent by a client
- `clock`: `int` for lamport clock, `[]int` for vector clock
2. Clients send a message every 5 seconds.
3. There's a 50-50 chance of dropping or forwarding received message for the server. (flag can be 0 or 1, if 0, forward message)
4. Clock is updated every time an action (i.e. receiving message, send message) is performed.
5. In `vector.go`, causality violation is checked every time a message is received (`clockIsGreaterThan(c1, c2)`) If the local vector clock of the receiving machine is more than the vector clock of the message, then a potential causality violation is detected.  

The `logs.txt` in each folder contains the sample outputs with 15 clients.

## Others
open issues
assumptions
