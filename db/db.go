package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec(
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			sessionId TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL
	);`)
	if err != nil {
		log.Fatal(err)
	}

}

func GetDB() *sql.DB {
	return DB
}
