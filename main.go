package main

import (
    "fmt"
    "net/http"
    "log"
    "database/sql"
    _ "modernc.org/sqlite"
)

type Document struct {
    ID int
    Title string
    Content string
}

func main() {
    db, err := sql.Open("sqlite", "./sqlite3.db") 
    if err != nil {
        log.Panic(err)
    }
    defer db.Close()
    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS documents (id INTEGER NOT NULL PRIMARY KEY, title TEXT, content TEXT);`)
    if err != nil {
        log.Panic(err)
    }
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    http.HandleFunc("/", handleIndex)
    http.HandleFunc("/view/", handleView)
    http.HandleFunc("/add", handleAdd)
    http.HandleFunc("/edit/", handleEdit)
    http.HandleFunc("/del/", handleDelete)
    fmt.Println("SERVER STARTING ON PORT 8080")
    http.ListenAndServe(":8080", nil)
}

