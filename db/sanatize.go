package db

import "strings"

const (
	// DB CRUD
	INS      = 0
	SEL      = 1
	UPD      = 2
	DEL      = 3
	CRUD_MAX = 6 // Longest length of sqlCRUD slice

)

// Prepares a string to create an SQL statement.
func prepString(crud int, table string, cols []string, condition string, valsLen int) string {
	retStr := make([]string, CRUD_MAX)
	copy(retStr, sqlCRUD[crud])

	colString := ""
	if crud != DEL {
		colString = strings.Join(cols, ", ")
	}

	switch crud {
	case INS:
		// ["INSERT INTO", "", "VALUES", ""]
		retStr[1] = table + " (" + colString + ")"
		retStr[3] = safeMarkers(len(cols), valsLen)

	case SEL:
		// ["SELECT (", "", ") FROM", "", "WHERE", ""]
		retStr[1] = colString
		retStr[3] = table
		retStr[5] = condition

	case UPD:
		// ["UPDATE", "", "SET (", "", ")WHERE", ""]
		retStr[1] = table
		retStr[3] = colString
		retStr[5] = condition

	case DEL:
		// ["DELETE FROM", "", "WHERE", ""]
		retStr[1] = table
		retStr[3] = condition
	}

	return strings.Join(retStr, "")
}

// Inserts safe characters, "?", into an sql statment.
//
// Used for INSERT commands
func safeMarkers(colsLen int, valsLen int) string {
	b := strings.Builder{}

	for row := 0; row < (valsLen / colsLen); row++ {
		b.WriteString("(")

		for i := 1; i < colsLen; i++ {
			b.WriteString("?,")
		}
		b.WriteString("?),")
	}

	return strings.TrimSuffix(b.String(), ",")
}

// Inserts safe characters, " = ?", following each column name.
//
// Used for UPDATE commands
func sanatizeCols(cols []string) []string {
	for i, v := range cols {
		cols[i] = v + " = ?"
	}
	return cols
}
