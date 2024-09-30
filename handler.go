package main

import (
    "fmt"
    "net/http"
)

func runRootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<h1>Home</h1><p>Welcome to the home page!</p>")
}

func runAddHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<h1>Add</h1><p>This is the add page!</p>")
}
