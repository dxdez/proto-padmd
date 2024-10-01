package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", runRootHandler)
    http.HandleFunc("/add", runAddHandler)

    fmt.Println("Server starting on PORT 8000")
    http.ListenAndServe(":8000", nil)
}

func runRootHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<h1>Home</h1><p>Welcome to the home page!</p>")
}

func runAddHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintln(w, "<h1>Add</h1><p>This is the add page!</p>")
}
