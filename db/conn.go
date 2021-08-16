package db

import (
	"database/sql"
	"log"
)

const (
	DB_DRIVER = "sqlite3"
	DB_PATH   = "/database/ittyurl.db"
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

func ConnectDB() (conn *DB) {
	db, err := sql.Open(DB_DRIVER, DB_PATH)
	if err != nil {
		dbFail(err.Error())
	}

	return &DB{db}
}
