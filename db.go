package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func openDB() error {
	db, err := sql.Open("sqlite3", "./sqlite3.db")
	if err != nil {
		return err
	}
	DB = db
	return nil
}

func closeDB() error {
	return DB.Close()
}

func setupDB() error {
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS items (id INTEGER NOT NULL PRIMARY KEY, title TEXT, content TEXT);`)
	if err != nil {
		return err
	}
	return nil
}
