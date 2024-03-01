package main

import (
    "log"

    "github.com/ProlificLabs/captrivia/server"
)

func main() {
    srv := server.NewServer()
    srv.StartServer()

    log.Println("Server has stopped")
}
