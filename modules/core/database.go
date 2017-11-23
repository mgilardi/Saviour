// Package core database checks database exists initializes the database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package core

import (
	"database/sql"
	// MySql
	_ "github.com/go-sql-driver/mysql"
)

// Database tyoe contains the sql access, options, logger, and the dsn for sql login
type Database struct {
	sql     *sql.DB
	options map[string]interface{}
	dsn     string
}

// InitDatabase initialize the database object and passes a pointer to the main loop
func InitDatabase() *Database {
	var db Database
	var err error
	var user, pass string
	DebugHandler.Sys("Starting", "DB")
	db.options = GetOptions("Core")
	if err != nil {
		DebugHandler.Err(err, "DB", 1)
	}
	if db.options["User"] == nil {
		DebugHandler.Err(err, "DB", 1)
	}
	if db.options["Pass"] == nil {
		DebugHandler.Err(err, "DB", 1)
	}
	user = db.options["User"].(string)
	pass = db.options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	DebugHandler.Sys("DSNLoaded", "DB")
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		DebugHandler.Err(err, "DB", 1)
	}
	err = db.sql.Ping()
	if err != nil {
		DebugHandler.Err(err, "DB", 1)
	}
	InitLogger(&db)
	InitCache(&db)
	db.CheckDB()
	return &db
}

// CheckDB checks if database exists and outputs tables that are found.
func (db *Database) CheckDB() {
	tables := make([]string, 0)
	rows, err := db.sql.Query(`SHOW TABLES`)
	if err != nil {
		LogHandler.Err(err, "DB")
		DebugHandler.Err(err, "DB", 1)
	}
	DebugHandler.Sys("AvaliableTables", "DB")
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			LogHandler.Err(err, "DB")
			DebugHandler.Err(err, "DB", 1)
		}
		DebugHandler.Sys(table, "DB")
		tables = append(tables, table)
	}
	db.createTables(tables)
	rows.Close()
}

// Runs SQL file if nothing exists in the database
func (db *Database) createTables(currentTables []string) {
	if len(currentTables) == 0 {
		// Load DB File
	}
}

// ResetIncrement keeps the increment tidy in the database
func (db *Database) ResetIncrement(tables ...string) {
	var maxIncrement int
	for _, table := range tables {
		db.sql.QueryRow(`SELECT MAX(UID) + 1 FROM ?`, table).Scan(&maxIncrement)
		db.sql.Exec(`ALTER TABLE ? AUTO_INCREMENT = ?`, table, maxIncrement)
	}
}
