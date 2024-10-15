package transport

import (
    "fmt"
    "net/http"
)

func RunRoutes() {
    fmt.Println("Setting up file server assets")
    fs := http.FileServer(http.Dir("./assets/static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))
    
    fmt.Println("Setting up routes")
    http.HandleFunc("/", runRootHandler)
    http.HandleFunc("/add", runAddHandler)

    fmt.Println("Server starting on PORT 8080")
    http.ListenAndServe(":8080", nil)
}
