package transport

import (
    "path/filepath"
    "html/template"
    "net/http"
    "log"
    "github.com/dylanxhernandez/proto-padmd/internal/models"
)

func RunRootHandler(w http.ResponseWriter, r *http.Request) {
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

func RunAddHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodPost {
        title := r.FormValue("title")
        if title == "" {
            return
        }
        _, err := models.InsertDocument(title)
        if err != nil {
            log.Printf("ERROR: %v", err)
        }
        // Render the form-reset template
        tmpl, err := template.ParseFiles("assets/templates/add.html")
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        err = tmpl.ExecuteTemplate(w, "page-content", nil)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
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
