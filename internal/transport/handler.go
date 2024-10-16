package transport

import (
    "html/template"
    "net/http"
    "log"
    "github.com/dylanxhernandez/proto-padmd/internal/models"
)

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
        }
        _, err := models.InsertDocument(title)
        if err != nil {
            log.Printf("ERROR: %v", err)
        }
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

