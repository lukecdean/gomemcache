package main

import (
    "gomemcache/server"
    "fmt"
)

func main() {
    s := server.GetInstance()
    fmt.Println("server host: " + server.Info(s))
    fmt.Println("running main")
} // main
