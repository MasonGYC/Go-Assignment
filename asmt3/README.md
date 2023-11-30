# implementations:
1. manager holds a priority queue, requests will be stored in queue and executed sequentially. Request with earlier timestamp or higher server id are prioritised.

# records
initially, every server possess one page, the page number of it is its own id

# TODO


go run vanilaIvy.go server.go record.go PriorityQueue.go  page.go message.go manager.go logger.go