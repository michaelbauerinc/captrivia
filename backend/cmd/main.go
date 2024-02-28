package main

import (
    "log"
    "github.com/ProlificLabs/captrivia/server" // Adjust the import path as necessary
)

func main() {
    server.StartServer()

    log.Println("Server has stopped")
}
