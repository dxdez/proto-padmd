package main

import (
    "fmt"
    "log"
    "github.com/dylanxhernandez/proto-padmd/internal/db"
    "github.com/dylanxhernandez/proto-padmd/internal/transport"
)

func main() {
    fmt.Println("Starting DB Connection")
    runOrError := db.OpenDB()
    if runOrError != nil {
    	log.Panic(runOrError)
    }
    defer db.CloseDB()
    runOrError = db.SetupDB()
    if runOrError != nil {
    	log.Panic(runOrError)
    }
    transport.RunRoutes()
}

