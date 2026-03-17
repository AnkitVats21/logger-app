package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./events.db?_pragma=foreign_keys(1)&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	queryUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT
	);`

	queryEvents := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		place TEXT,
		timestamp DATETIME,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err := DB.Exec(queryUsers)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = DB.Exec(queryEvents)
	if err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}
}
