package main

import (
    "html/template"
    "fmt"
    "net/http"
    "regexp"
    "log"
    "database/sql"
    "github.com/gomarkdown/markdown"
    "github.com/gomarkdown/markdown/html"
    "github.com/gomarkdown/markdown/parser"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
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
        if len(currentDocument.Content) >= 100 {
    	currentDocument.Content = currentDocument.Content[:100]
        } 
        documents = append(documents, currentDocument)
    }
    tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/content.html"))
    err = tmpl.ExecuteTemplate(w, "base", map[string]any{"Documents": documents})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handleView(w http.ResponseWriter, r *http.Request) {
    idEditRegex := regexp.MustCompile(`^/view/([0-9]+)$`)
    matches := idEditRegex.FindStringSubmatch(r.URL.Path)
    if len(matches) != 2 {
        http.NotFound(w, r)
        return
    }
    id := matches[1]     
    var document Document
    err := db.QueryRow("SELECT id, title, content FROM documents WHERE id = ?", id).Scan(&document.ID, &document.Title, &document.Content)            
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    
    extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
    p := parser.NewWithExtensions(extensions)
    returnDoc := p.Parse([]byte(document.Content))
    
    htmlFlags := html.CommonFlags | html.HrefTargetBlank
    opts := html.RendererOptions{Flags: htmlFlags}
    renderer := html.NewRenderer(opts)
    
    htmlContent := markdown.Render(returnDoc, renderer)
    htmlContentStr := string(htmlContent)        
    fmt.Println(htmlContentStr)
    tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/content_markdown.html"))
    err = tmpl.ExecuteTemplate(w, "base", template.HTML(htmlContentStr))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        title := r.FormValue("title")
        content := r.FormValue("content")
        if title == "" {
    	    return
        }
        result, err := db.Exec("INSERT INTO documents (title, content) VALUES (?, ?)", title, content)
        if err != nil {
    	    log.Printf("ERROR: %v", err)
        }
        id, err := result.LastInsertId()
        if err != nil {
    	    log.Printf("ERROR: %v", err)
        }
        tmpl := template.Must(template.ParseFiles("templates/form_submitted.html"))
        err = tmpl.ExecuteTemplate(w, "content", map[string]any{"IdRef": id})
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
}

func handleEdit(w http.ResponseWriter, r *http.Request) {
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
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
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
}

