package web

import (
	"database/sql"
	"log"

	"github.com/vinm0/ittyurl/data"
	"github.com/vinm0/ittyurl/db"
)

const (
	USRS_TABLE = "usrs"
)

func DBAuth(email string) (usr *data.User, auth bool) {
	conn := db.ConnectDB()
	defer conn.Close()

	rows, err := conn.Select(USRS_TABLE, usrColsAll(), "email = ?", usr.Email)
	if err != nil {
		log.Println(err.Error())
	}

	dbUsr := loadUser(rows)

	return dbUsr, (dbUsr != nil)
}

func loadUser(rows *sql.Rows) *data.User {
	defer rows.Close()

	var usr *data.User = nil
	if rows.Next() {
		usr = &data.User{}
		rows.Scan(usr.UserID, usr.Fname, usr.Lname, usr.Email,
			usr.Usrname, usr.Joindate, usr.AcctID, usr.Pwd)
	}

	return usr
}

func usrColsAll() (cols []string) {
	return usrColsPwd(true)
}

// func usrCols() (cols []string) {
// 	return usrColsPwd(false)
// }

func usrColsPwd(pwd bool) (cols []string) {
	cols = []string{"usr_id", "fname", "lname", "email",
		"usrname", "joindate", "acct_id",
	}

	if pwd {
		cols = append(cols, "pwd")
	}

	return cols
}
