package main

import (
    "fmt"
    "html/template"
    "net/http"
    "log"
    "path/filepath"
)

func main() {
    fmt.Println("Starting DB Connection")

    runOrError := openDB()
    if runOrError != nil {
    	log.Panic(runOrError)
    }
    defer closeDB()
    runOrError = setupDB()
    if runOrError != nil {
    	log.Panic(runOrError)
    }

    http.HandleFunc("/", runRootHandler)
    http.HandleFunc("/add", runAddHandler)

    fmt.Println("Server starting on PORT 8080")
    http.ListenAndServe(":8080", nil)
}

func runRootHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "home.html", nil)
}

func runAddHandler(w http.ResponseWriter, r *http.Request) {
    renderTemplate(w, "add.html", nil)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    templatesDir := "templates"
    files := []string{
        filepath.Join(templatesDir, tmpl),
        filepath.Join(templatesDir, "layout.html"),
    }
    templateServe, err := template.ParseFiles(files...)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    err = templateServe.ExecuteTemplate(w, "layout", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
