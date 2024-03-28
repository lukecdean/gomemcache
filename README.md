## Basic implementation of memcache in Go to learn the language
memcache is an in memory key-value store that is interfaced via TCP
Start server:   go run main.go s
Start a client: go run main.go c
Client commands:
  set <key> <value>
  get <key>
  exit

### Goals
- general familiarity with Go
- sockets in Go
- goroutines

### What I actually learned
- general familiarity with how to program with Go
  - syntax, defer (really cool), 'class' (package) structure, building and running, etc
- object oriented programming in Go (it's not like Java)
- using goroutines
- using net sockets in Go
- parsing tcp messages
- using mutexes in Go
- working with strings in Go
- singleton design pattern (for the server)
