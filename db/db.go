package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./edu_resources.db")
	if err != nil {
		log.Fatal(err)
	}

	sqlStmt := `
	CREATE TABLE IF NOT EXISTS resources (id INTEGER PRIMARY KEY, title TEXT, category TEXT, description TEXT, url TEXT, date_added DATETIME, resource_type TEXT, completion_time TEXT);
	`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sqlStmt)
	}
}