package client

import (
    "fmt"
    "net"
)

type client struct {
    serverHost  string // "localhost"
    serverPort  string // "3333"
    serverType  string // "tcp"
} // struct

func New() client {
    c := client {"localhost", "3333", "tcp"}
    return c
} // New()

func Connect(c client) {
    connection, err := net.Dial(c.serverType, c.serverHost + ":" + c.serverPort)
    if err != nil {
        panic(err)
    } // if err not nil

    _, err = connection.Write([]byte("Hello, server"))
    buffer := make([]byte, 1024)
    mLen, err := connection.Read(buffer)
    if err != nil {
        fmt.Println("err reading: ", err.Error())
    } // if err

    fmt.Println("recieved: ", string(buffer[:mLen]))
    defer connection.Close()
} // Connect()
