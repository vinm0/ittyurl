package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/vinm0/ittyurl/db"
)

const (
	// Database Table Names
	USRS_TABLE = "usrs"
)

// FindUsrDB obtains client data from the database (password included).
// Returns usr populated with client data if the clien is found.
// Otherwise, returns nil.
// FindUsrDB is intended for obtaining data for client login.
// To obtain data for a user other than the current client, please use UsrDB.
func FindUsr(email string) (usr *User, found bool) {
	conn := db.ConnectDB()
	defer conn.Close()

	row, err := conn.Select(USRS_TABLE, usrColsAll(), "email = ?", usr.Email)
	if err != nil {
		log.Println(err.Error())
	}

	dbUsr := loadUser(row, true)

	return dbUsr, (dbUsr != nil)
}

func FindUrl(path string) (url *Url, found bool) {
	conn := db.ConnectDB()
	defer conn.Close()

	row, err := conn.Select(URLS_TABLE, urlUsrCols(), "alt = ?", path)
	if err != nil {
		log.Println(err.Error())
	}

	dbUsr := loadUrl(row)

	return dbUsr, (dbUsr != nil)
}

// Inserts the User data into the database.
func createUser(usr *User) error {
	return insert(usr)
}

// Inserts the Url data into the database.
func insertUrl(url *Url) error {
	return insert(url)
}

// UrlBySource returns the Url instance with the corresponding src url address.
//
// Returns nil, if no Url is found.
func UrlBySource(src string) *Url {
	conn := db.ConnectDB()
	defer conn.Close()

	rows, _ := conn.Select(URLS_TABLE, urlUsrCols(), "source = ?", src)
	fmt.Println("Source url:", src)

	return loadUrl(rows)
}

// Increment's the visit count for the existing path.
func IncrementVisits(path string) {
	conn := db.ConnectDB()
	defer conn.Close()

	cols := db.Cols("viscount = (viscount + 1)")
	conn.Update("untracked_visits", cols, "alt = ?", path)
}

func CountUrlsByUserID(id int) int {
	return countUrls(id)
}

func CountUrlsByIP(ip string) int {
	return countUrls(ip)
}

// loadUser populates a new User instance with fields from an SQL row.
// Use pwd to indicate whether the User password is included in the query.
// Returns a pointer to the populated instance if the row contains fields.
// Otherwise, returns nil.
func loadUser(row *sql.Rows, pwd bool) (usr *User) {
	defer row.Close()

	if row.Next() {
		usr = &User{}

		if pwd {
			row.Scan(usr.UserID, usr.Fname, usr.Lname, usr.Email,
				usr.Usrname, usr.Joindate, usr.AcctID, usr.Pwd)
		} else {
			row.Scan(usr.UserID, usr.Fname, usr.Lname, usr.Email,
				usr.Usrname, usr.Joindate, usr.AcctID)
		}
	}

	return usr
}

// Populates a new Url instance with fields from an SQL row.
// Returns a pointer to the populated instance if the row contains fields.
// Otherwise, returns nil.
func loadUrl(row *sql.Rows) (url *Url) {
	defer row.Close()

	if row.Next() {
		url = &Url{}
		usr := &User{}
		urlDate := ""
		usrDate := ""
		row.Scan(&url.Path, &url.Source, &urlDate, &usr.UserID,
			&url.CreatorIP, &usr.UserID, &usr.Fname, &usr.Lname, &usr.Email,
			&usr.Usrname, &usrDate, &usr.AcctID,
		)

		usr.Joindate, _ = time.Parse(TIME_FORMAT, usrDate)
		url.DateCreated, _ = time.Parse(TIME_FORMAT, urlDate)
		url.Owner = usr
	}

	return url
}

// usrColsAll returns a slice of all the column names for the user table
// (including password).
// To obtain the column names without the password column,
// use usrCols instead.
func usrColsAll() (cols []string) {
	return usrColsPwd(true)
}

// usrCols returns a slice of the column names for the user table
// (excluding password).
// To obtain the column names including the password column,
// use usrColsAll instead.
func usrCols() (cols []string) {
	return usrColsPwd(false)
}

// usrColsPwd returns a slice of the column names for the user table
// The password column is included if pwd is set to true.
func usrColsPwd(pwd bool) (cols []string) {
	cols = []string{"usr_id", "fname", "lname", "email",
		"usrname", "joindate", "acct_id",
	}

	if pwd {
		cols = append(cols, "pwd")
	}

	return cols
}

// Returns the column names for a JOINed usrs and urls table.
//
// ** Used for retreiving data for a Url instance (including Owner) **
func urlUsrCols() (cols []string) {
	return append(urlCols(), usrCols()...)
}

// Returns the column names for the urls table.
func urlCols() (cols []string) {
	return []string{
		"alt",
		"source",
		"datecreated",
		"owner_id",
		"creatorip",
	}
}

// Returns the vals of the obj of a predetermined type.
//
// *User | *Url
func Vals(obj interface{}) []interface{} {
	if usr, ok := obj.(*User); ok {
		return db.Vals(usr.UserID, usr.Fname, usr.Lname, usr.Email,
			usr.Usrname, usr.Joindate.Format(TIME_FORMAT), usr.AcctID, usr.Pwd)
	}

	if url, ok := obj.(*Url); ok {
		return db.Vals(url.Path, url.Source, url.DateCreated.Format(TIME_FORMAT),
			url.Owner.UserID, url.CreatorIP)
	}

	return nil
}

// Inserts data into the database of objects of a predetermined type
//
// *User | *Url
func insert(data interface{}) (err error) {
	conn := db.ConnectDB()
	defer conn.Close()

	switch data.(type) {
	case *User:
		_, err = conn.Insert(USRS_TABLE, usrCols(), Vals(data)...)
	case *Url:
		_, err = conn.Insert("urls", urlCols(), Vals(data)...)
	}

	return err
}

// Returns the row count of Urls given a specified condition
func countUrls(id interface{}) int {
	conn := db.ConnectDB()
	defer conn.Close()

	cols := db.Cols("COUNT(*)")

	var rows *sql.Rows
	switch id.(type) {
	case int: // UserID
		rows, _ = conn.Select("urls", cols, "owner_id = ?", id)
	case string: // IP
		rows, _ = conn.Select("urls", cols, "creatorip = ?", id)
	}
	defer rows.Close()

	var count int
	if rows.Next() {
		rows.Scan(&count)
	}

	return count
}
