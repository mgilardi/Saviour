package core

import (
	"testing"
)

var db *Database
var uid int

func TestDatabase_InitDatabase(t *testing.T) {
	stmt, err := DBHandler.sql.Prepare(`DELETE FROM logger WHERE type = ?`)
	deleteFromDB := make([]interface{}, 0)
	deleteFromDB = append(deleteFromDB, stmt, "Error")
	if err != nil {
		Logger(err.Error(), "DB", ERROR)
	}
	DBHandler.Exec(deleteFromDB)
}
