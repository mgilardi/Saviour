package core

import (
	"testing"
)

var db *Database
var uid int

func TestDatabase_InitDatabase(t *testing.T) {
	InitDebug(true)
	InitCron()
	db = InitDatabase()
}
