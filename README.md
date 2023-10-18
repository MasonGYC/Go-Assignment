# Go-Assignment
50.041 Distributed Systems and Computing Go Assignment 

`go mod init go_asmt/asmt1`
Q1:
`go run ./lamport/lamport.go`
`go run ./vector/vector.go`

Q2: 
`go run .\bully.go .\common.go .\message.go .\node.go .\node_listen.go .\logger.go`

todo:
1. lamport total order, present order
2. q1 register pattern
3. q1 function func (a * server), don't use go func () inside a func
4. change to printf
5. timestamp in message change to logical clock
6. mutex for n.state
7. break -> return 
8. write tests (order of events)