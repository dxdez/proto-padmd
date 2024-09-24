package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintln(w, "<h1>Home</h1><p>Welcome to the home page!</p>")
    })

    http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintln(w, "<h1>Test</h1><p>This is the test page!</p>")
    })

    fmt.Println("Server starting on PORT 8000")
    http.ListenAndServe(":8000", nil)
}
