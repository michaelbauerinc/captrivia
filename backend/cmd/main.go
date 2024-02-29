package main

import (
    "log"
    "github.com/ProlificLabs/captrivia/server"
)

func main() {
    server.StartServer()

    log.Println("Server has stopped")
}
