package db

import (
	"database/sql"
	"fmt"
)

func Cols(columns ...string) []string {
	return columns
}

func Vals(values ...interface{}) []interface{} {
	return values
}

func (db *DB) Select(table string, cols []string, condition string, vals ...interface{}) (*sql.Rows, error) {
	prepStr := prepString(SEL, table, cols, condition, len(vals))

	stmt, err := db.Prepare(prepStr)
	if err != nil {
		dbFail(err.Error())
	}
	defer stmt.Close()

	rows, err := stmt.Query(vals...)
	if err != nil {
		dbFail(err.Error())
	}

	return rows, err
}

func (db *DB) Insert(table string, cols []string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(INS, table, cols, vals, "", nil)
}

func (db *DB) InsertMany(table string, cols []string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(INS, table, cols, vals, "", nil)
}

func (db *DB) Update(table string, cols []string, condition string, vals ...interface{}) (sql.Result, error) {
	cols = sanatizeCols(cols)
	return db.sqlGeneric(UPD, table, cols, vals, condition, nil)
}

func (db *DB) Delete(table string, conditions string, vals ...interface{}) (sql.Result, error) {

	return db.sqlGeneric(DEL, table, nil, vals, conditions, nil)
}

func (db *DB) sqlGeneric(crud int, table string, cols []string,
	vals []interface{}, condition string, rows *sql.Rows) (sql.Result, error) {

	prepStr := prepString(crud, table, cols, condition, len(vals))
	fmt.Println(prepStr)

	stmt, err := db.Prepare(prepStr)
	if err != nil {
		dbFail(err.Error())
	}

	res, err := stmt.Exec(vals...)
	if err != nil {
		dbFail(err.Error())
	}

	return res, err
}
