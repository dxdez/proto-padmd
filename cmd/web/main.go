package main

import (
    "fmt"
    "html/template"
    "net/http"
    "math/rand"
    "time"
    "log"
    "path/filepath"
    "github.com/dylanxhernandez/proto-padmd/internal/db"
    "github.com/dylanxhernandez/proto-padmd/internal/models"
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

    http.HandleFunc("/", runRootHandler)
    http.HandleFunc("/add", runAddHandler)

    fmt.Println("Server starting on PORT 8080")
    http.ListenAndServe(":8080", nil)
}

func runRootHandler(w http.ResponseWriter, r *http.Request) {
    documents, error := models.GetAllDocuments() 
    if error != nil {
        log.Printf("ERROR: %v", error)
        return
    }
    data := models.DocumentLists {
        Documents: documents,
    }
    renderTemplate(w, "index.html", data)
}

func runAddHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
	title := r.FormValue("title")
	if title == "" {
            return
	} else {
            fmt.Println(fmt.Sprintf("The title is: %s", title))
        }
        rand.Seed(time.Now().UnixNano())
        randomNumber := rand.Intn(100)
        randomNumberString := fmt.Sprintf("POST Title %d", randomNumber)
        fmt.Println(randomNumberString)

        _, error := models.InsertDocument(randomNumberString)
        if error != nil {
            log.Printf("ERROR: %v", error)
        }

        http.Redirect(w, r, "/add", http.StatusSeeOther)
        return
    } else {
        renderTemplate(w, "add.html", nil)
    }
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    templatesDir := "assets/templates"
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
