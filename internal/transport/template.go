package transport

import (
    "path/filepath"
    "html/template"
    "net/http"
)

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
