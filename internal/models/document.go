package models

import (
    "github.com/dylanxhernandez/proto-padmd/internal/db"
)

type Document struct {
    ID int
    Title string
    Content string
}

type DocumentLists struct {
    Documents []Document
}

func GetAllDocuments() ([]Document, error) {
    var documentList []Document

    documentRows, err := db.DB.Query("SELECT id, title, content FROM documents")
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

func InsertDocument(title string) (Document, error) {
    var id int
    err := db.DB.QueryRow("INSERT INTO documents (title, content) VALUES (?, 'This is sample content') RETURNING id", title).Scan(&id)
    if err != nil {
        return Document{}, err
    }
    document := Document{ID: id, Title: title, Content: "This is sample content"}
    return document, nil
}
