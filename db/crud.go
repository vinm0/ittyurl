package db

import (
	"database/sql"
	"fmt"
)

// Converts a list of column names to a slice.
func Cols(columns ...string) []string {
	return columns
}

// Converts a list of values to a slice.
func Vals(values ...interface{}) []interface{} {
	return values
}

// Executes a SELECT command on the database.
// Returns the rows object containing the results of the query.
func (db *DB) Select(table string, cols []string,
	condition string, vals ...interface{}) (*sql.Rows, error) {

	stmt := db.sqlStatement(SEL, table, cols, condition, len(vals))
	defer stmt.Close()

	rows, err := stmt.Query(vals...)
	if err != nil {
		dbFail(err.Error())
	}

	return rows, err
}

// Executes an INSERT command on the database.
// Returns the results of the statement.
func (db *DB) Insert(table string, cols []string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(INS, table, cols, vals, "")
}

// Executes an INSERT command on the database.
// Returns the results of the statement.
func (db *DB) InsertMany(table string, cols []string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(INS, table, cols, vals, "")
}

// Executes an UPDATE command on the database.
// Returns the results of the statement.
func (db *DB) Update(table string, cols []string, condition string, vals ...interface{}) (sql.Result, error) {
	cols = sanatizeCols(cols)
	return db.sqlGeneric(UPD, table, cols, vals, condition)
}

// Executes an DELETE command on the database.
// Returns the results of the statement.
func (db *DB) Delete(table string, conditions string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(DEL, table, nil, vals, conditions)
}

// Called by CRUD functions to execute commands on the database
// Returns the results of the statement.
func (db *DB) sqlGeneric(crud int, table string, cols []string,
	vals []interface{}, condition string) (sql.Result, error) {

	stmt := db.sqlStatement(crud, table, cols, condition, len(vals))
	defer stmt.Close()

	res, err := stmt.Exec(vals...)
	if err != nil {
		dbFail(err.Error())
	}

	return res, err
}

// Creates a prepared SQL statment to issue database commands.
//
// Call "defer stmt.Close() after creating the statment"
func (db *DB) sqlStatement(crud int, table string, cols []string,
	condition string, valsLen int) (stmt *sql.Stmt) {

	prepStr := prepString(crud, table, cols, condition, valsLen)
	fmt.Println(prepStr)

	stmt, err := db.Prepare(prepStr)
	if err != nil {
		dbFail(err.Error())
	}

	return stmt
}
