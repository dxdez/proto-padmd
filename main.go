package main

import (
        "html/template"
        "fmt"
        "net/http"
        "regexp"
        "log"
        "database/sql"
        _ "modernc.org/sqlite"
       )

var DB *sql.DB

type Document struct {
    ID int
    Title string
    Content string
}

func main() {
    fmt.Println("STARTING DB CONNECTION")
    db, err := sql.Open("sqlite", "./sqlite3.db") 
    if err != nil {
        log.Panic(err)
    }
    DB = db
    defer DB.Close()
    _, err = DB.Exec(`CREATE TABLE IF NOT EXISTS documents (id INTEGER NOT NULL PRIMARY KEY, title TEXT, content TEXT);`)
    if err != nil {
        log.Panic(err)
    }

    fmt.Println("SETTING UP STATIC ASSETS")
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
    
    fmt.Println("SETTING UP ROUTES")
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        documents, error := getAllDocuments() 
        if error != nil {
            log.Printf("ERROR: %v", error)
            return
        }
        tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/content.html"))
        err := tmpl.ExecuteTemplate(w, "base", map[string]any{"Documents": documents})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    http.HandleFunc("/add", func (w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            title := r.FormValue("title")
            if title == "" {
                return
            }
            _, err := insertDocument(title)
            if err != nil {
                log.Printf("ERROR: %v", err)
            }
            tmpl := template.Must(template.ParseFiles("templates/form_submitted.html"))
            err = tmpl.ExecuteTemplate(w, "content", nil)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        } else {
            tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/form.html"))
            err := tmpl.ExecuteTemplate(w, "base", nil)
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        }
    })
    http.HandleFunc("/del/", func (w http.ResponseWriter, r *http.Request) {
        idDelRegex := regexp.MustCompile(`^/del/([0-9]+)$`)
        matches := idDelRegex.FindStringSubmatch(r.URL.Path)
        if len(matches) != 2 {
            http.NotFound(w, r)
            return
        }
        
        id := matches[1]     
        _, err := DB.Exec("DELETE FROM documents WHERE id = (?)", id)
        if err != nil {
            log.Printf("ERROR: %v", err)
        }
        http.Redirect(w, r, "/", http.StatusFound)
    })
    
    fmt.Println("SERVER STARTING ON PORT 8080")
    http.ListenAndServe(":8080", nil)
}

func getAllDocuments() ([]Document, error) {
    var documentList []Document

    documentRows, err := DB.Query("SELECT id, title, content FROM documents")
    if err != nil {
        return nil, err
    }
    defer documentRows.Close()

    for documentRows.Next() {
        currentDocument := Document{}
        err := documentRows.Scan(&currentDocument.ID, &currentDocument.Title, &currentDocument.Content)
        if err != nil {
            return []Document{}, err
        }
        documentList = append(documentList, currentDocument)
    }

    return documentList, nil
}

func insertDocument(title string) (Document, error) {
    var id int
    err := DB.QueryRow("INSERT INTO documents (title, content) VALUES (?, 'This is sample content') RETURNING id", title).Scan(&id)
    if err != nil {
        return Document{}, err
    }
    document := Document{ID: id, Title: title, Content: "This is sample content"}
    return document, nil
}

