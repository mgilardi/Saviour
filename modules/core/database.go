// Package core database checks database exists initializes the database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package core

import (
	"database/sql"
	// MySql
	_ "github.com/go-sql-driver/mysql"
)

// DBHandler Database Global Variable
var DBHandler *Database

// Database tyoe contains the sql access, options, logger, and the dsn for sql login
type Database struct {
	sql *sql.DB
	dsn string
}

// InitDatabase initialize the database object and passes a pointer to the main loop
func InitDatabase() {
	var db Database
	var err error
	var user, pass string
	Sys("Starting", "DB")
	options := OptionsHandler.GetOptions("Core")
	if err != nil {
		Error(err, "DB")
	}
	if options["User"] == nil {
		Error(err, "Cache")
	}
	if options["Pass"] == nil {
		Error(err, "Cache")
	}
	user = options["User"].(string)
	pass = options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	Sys("DSNLoaded", "DB")
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		Error(err, "Cache")
	}
	err = db.sql.Ping()
	if err != nil {
		Error(err, "Cache")
	}

	db.checkDB()
	DBHandler = &db
}

// CheckDB checks if database exists and outputs tables that are found.
func (db *Database) checkDB() {
	tables := make([]string, 0)
	rows, err := db.sql.Query(`SHOW TABLES`)
	if err != nil {
		Error(err, "Cache")
	}
	Sys("AvaliableTables", "DB")
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			Error(err, "Cache")
		}
		Sys(table, "DB")
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
