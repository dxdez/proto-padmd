package main

import (
    "fmt"
    "net/http"
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

    fs := http.FileServer(http.Dir("./assets/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    http.HandleFunc("/", transport.RunRootHandler)
    http.HandleFunc("/add", transport.RunAddHandler)

    fmt.Println("Server starting on PORT 8080")
    http.ListenAndServe(":8080", nil)
}

