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
	Logger("Starting", "DB", MSG)
	options := OptionsHandler.GetOptions("Core")
	if options["User"] == nil {
		Logger(err.Error(), "DB", FATAL)
	}
	if options["Pass"] == nil {
		Logger(err.Error(), "DB", FATAL)
	}
	user = options["User"].(string)
	pass = options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	Logger("DSNLoaded", "DB", MSG)
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		Logger(err.Error(), "DB", FATAL)
	}
	err = db.sql.Ping()
	if err != nil {
		Logger(err.Error(), "DB", FATAL)
	}

	db.checkDB()
	DBHandler = &db
}

// CheckDB checks if database exists and outputs tables that are found.
func (db *Database) checkDB() {
	tables := make([]string, 0)
	rows, err := db.sql.Query(`SHOW TABLES`)
	if err != nil {
		Logger(err.Error(), "DB", ERROR)
	}
	Logger("AvaliableTables", "DB", MSG)
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			Logger(err.Error(), "Cache", ERROR)
		}
		Logger(table, "DB", MSG)
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
