package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	DB_DRIVER = "sqlite3"
	DB_PATH   = "./database/ittyurl.db"
)

var (
	dbFail func(...interface{}) = log.Println

	// SQL command templates
	sqlCRUD = [][]string{
		{"INSERT INTO ", "", " VALUES ", ""},
		{"SELECT ", "", " FROM ", "", " WHERE ", ""},
		{"UPDATE ", "", " SET ", "", " WHERE ", ""},
		{"DELETE FROM ", "", " WHERE ", ""},
	}
)

type DB struct{ *sql.DB }

// Creates a connection to the database.
//
// Please use "defer conn.Close()" to close the connection.
func ConnectDB() (conn *DB) {
	db, err := sql.Open(DB_DRIVER, DB_PATH)
	if err != nil {
		dbFail(err.Error())
	}

	return &DB{db}
}
