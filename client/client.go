package client

import (
    "bufio"
    "fmt"
    "net"
    "os"
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
    defer connection.Close()

    reader := bufio.NewReader(os.Stdin)

    for {
        // take input
        fmt.Print("-> ")
        text, _ := reader.ReadString('\n')
        //text = strings.Replace(text, "\n", "", -1)

        if text == "exit\n" {
            os.Exit(0)
        } // if

        _, err = connection.Write([]byte(text))

        buffer := make([]byte, 1024)
        mLen, err := connection.Read(buffer)
        if err != nil {
            fmt.Println("err reading: ", err.Error())
        } // if err

        fmt.Println("recieved: ", string(buffer[:mLen]))
    } // for 
} // Connect()
