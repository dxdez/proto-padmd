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
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        var documents []Document
        documentRows, err := db.Query("SELECT id, title, content FROM documents")
        if err != nil {
            log.Panic(err)
        }
        defer documentRows.Close()
        for documentRows.Next() {
            currentDocument := Document{}
            err = documentRows.Scan(&currentDocument.ID, &currentDocument.Title, &currentDocument.Content)
            if err != nil {
               log.Panic(err)
            }
            if len(currentDocument.Content) >= 75 {
                currentDocument.Content = currentDocument.Content[:75]
            } 
            documents = append(documents, currentDocument)
        }
        tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/content.html"))
        err = tmpl.ExecuteTemplate(w, "base", map[string]any{"Documents": documents})
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    })
    http.HandleFunc("/add", func (w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            title := r.FormValue("title")
            content := r.FormValue("content")
            if title == "" {
                return
            }
            _, err := db.Exec("INSERT INTO documents (title, content) VALUES (?, ?)", title, content)
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
    http.HandleFunc("/edit/", func (w http.ResponseWriter, r *http.Request) {
        idEditRegex := regexp.MustCompile(`^/edit/([0-9]+)$`)
        matches := idEditRegex.FindStringSubmatch(r.URL.Path)
        if len(matches) != 2 {
            http.NotFound(w, r)
            return
        }
        id := matches[1]     
        if r.Method == http.MethodPost {
            title := r.FormValue("title")
            content := r.FormValue("content")
            if title == "" {
                return
            }
            _, err := db.Exec("UPDATE documents SET title = (?), content = (?) WHERE id = (?)", title, content, id)
            if err != nil {
                log.Printf("ERROR: %v", err)
            }
            tmpl := template.Must(template.ParseFiles("templates/form_submitted.html"))
            err = tmpl.ExecuteTemplate(w, "content", map[string]any{"IdRef": id})
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
        } else {
            var document Document
            err := db.QueryRow("SELECT id, title, content FROM documents WHERE id = ?", id).Scan(&document.ID, &document.Title, &document.Content)            
            if err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
            }
            tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/form.html"))
            err = tmpl.ExecuteTemplate(w, "base", map[string]any{"Editing": true, "IdRef": id, "TitleRef": document.Title, "ContentRef": document.Content })
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
        _, err := db.Exec("DELETE FROM documents WHERE id = (?)", id)
        if err != nil {
            log.Printf("ERROR: %v", err)
        }
        http.Redirect(w, r, "/", http.StatusFound)
    })
    fmt.Println("SERVER STARTING ON PORT 8080")
    http.ListenAndServe(":8080", nil)
}

