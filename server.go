package main
import (
    "fmt"
    "net"
    "os"
)
const (
    SERVER_HOST = "localhost"
    SERVER_PORT = "3333"
    SERVER_TYPE = "tcp"
)
func main() {
    fmt.Println("server booting")
    server, err := net.Listen(
        SERVER_TYPE, SERVER_HOST + ":" + SERVER_PORT
    )
    if err != nil {
        fmt.Println("err listening: ", err.Error())
        os.Exit(1)
    } // if err
    defer server.Close()
    fmt.Println(
        "server listening on "
        + SERVER_HOST + ":" + SERVER_PORT
    )
    for {
        connection, err := server.Accept()
        if err != nil {
            fmt.Println("err accepting: ", err.Error())
            os.Exit(1)
        } // if
        fmt.Println("client connected: " + connection.RemoteAddr())
        go processClient(connection)
    } // for
} // main()
