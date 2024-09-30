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
