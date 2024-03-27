package server

import (
    "fmt"
    "net"
    "os"
)

var instance    *server
var data        map[string]string

type server struct {
    serverHost  string // "localhost"
    serverPort  string // "3333"
    serverType  string // "tcp"
    running     bool
}

// singleton pattern
func GetInstance() *server {
    if instance == nil {
        instance = &server{
            serverHost: "localhost",
            serverPort: "3333",
            serverType: "tcp",
            running:    false,
        }
    }
    return instance
}

// boot the server and make it start listening
func Boot(s *server) {
    if s.running {
        fmt.Println("server already running")
        return
    } // if running
    fmt.Println("server booting")

    // create key value map
    data = make(map[string]string)

    server, err := net.Listen(s.serverType, s.serverHost + ":" + s.serverPort)
    if err != nil {
        fmt.Println("err listening: ", err.Error())
        os.Exit(1)
    } // if err
    fmt.Println("server listening on " + s.serverHost + ":" + s.serverPort)

    defer server.Close()
    s.running = true
    defer cleanUp(s)

    // listen for clients
    for {
        connection, err := server.Accept()
        if err != nil {
            fmt.Println("err accepting: ", err.Error())
            os.Exit(1)
        } // if
        fmt.Println("client connected: " + connection.RemoteAddr().String())
        go serveClient(connection)
    } // for
} // main()

// set members to non-running values
func cleanUp(s *server) {
    s.running = false
    data = nil
} // cleanUp()

// complete requests from a client
func serveClient(connection net.Conn) {
} // serveClient()

// honestly a func just for testing
func Info(s *server) (string) {
    return s.serverHost
} // info()
