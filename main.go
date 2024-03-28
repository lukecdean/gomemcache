package main

import (
    "gomemcache/server"
    "gomemcache/client"
    "fmt"
    "os"
)

func main() {
    fmt.Println("running main")

    switch os.Args[1] {
    case "s":
        fallthrough
    case "server":
        runServer()
    case "c":
        fallthrough
    case "client":
        runClient()
    default:
        fmt.Println("client or server?")
    } // switch

} // main

func runServer() {
    s := server.GetInstance()
    server.Boot(s)
} // runServer()

func runClient() {
    c := client.New()
    client.Connect(c)
} // runClient()
