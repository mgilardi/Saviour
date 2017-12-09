// Package core database checks database exists initializes the database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package core

import (
	"database/sql"
	// MySql
	_ "github.com/go-sql-driver/mysql"
)

const (
	// PACKAGE is the package name constant
	PACKAGE = "Core"
	// MODULEDB is the module name constant
	MODULEDB = "Database"
)

// DBHandler Database Global Variable
var DBHandler *Database

// Database tyoe contains the sql access, options, logger, and the dsn for sql login
type Database struct {
	sql *sql.DB
	dsn string
}

// @TODO Put wrapper around SQL statements for transactions

// InitDatabase initialize the database object and passes a pointer to the main loop
func InitDatabase() {
	var db Database
	var err error
	var user, pass string
	Logger("Starting", PACKAGE+"."+MODULEDB+".InitDatabase", MSG)
	options := OptionsHandler.GetOptions("Core")
	if options["User"] == nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".InitDatabase", FATAL)
	}
	if options["Pass"] == nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".InitDatabase", FATAL)
	}
	user = options["User"].(string)
	pass = options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	Logger("DSNLoaded", PACKAGE+"."+MODULEDB+".InitDatabase", MSG)
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".InitDatabase", FATAL)
	}
	err = db.sql.Ping()
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".InitDatabase", FATAL)
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
	Logger("AvaliableTables", PACKAGE+"."+MODULEDB+".CheckDB", MSG)
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

// Query executes a query call to the database
func (db *Database) Query(exec string, args ...interface{}) *sql.Row {
	var row *sql.Row
	row = db.sql.QueryRow(exec, args...)
	return row
}

// QueryRows returns a query of multipule rows
func (db *Database) QueryRows(exec string, args ...interface{}) (bool, *sql.Rows) {
	var rows *sql.Rows
	exists := true
	rows, err := db.sql.Query(exec, args...)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".QueryRows", ERROR)
		exists = false
	}
	return exists, rows
}

// SetupExec sets up the SQL exec command and input variables to be executed
func (db *Database) SetupExec(exec string, args ...interface{}) []interface{} {
	setupOutput := make([]interface{}, 0)
	stmt, err := db.sql.Prepare(exec)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEDB+".SetupExec", ERROR)
	}
	setupOutput = append(setupOutput, stmt)
	for _, arg := range args {
		setupOutput = append(setupOutput, arg)
	}
	return setupOutput
}

// Exec will execute the SQL commands
func (db *Database) Exec(querys ...[]interface{}) {
	tx, err := db.sql.Begin()
	if len(querys) == 0 {
		Logger("QueryWrap::NoElements", PACKAGE+"."+MODULEDB+".Exec", ERROR)
	} else {
		for _, query := range querys {
			stmt := query[0].(*sql.Stmt)
			query = append(query[:0], query[0+1:]...)
			_, err = tx.Stmt(stmt).Exec(query...)
		}
		if err != nil {
			Logger(err.Error(), PACKAGE+"."+MODULEDB+".Exec", ERROR)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
}

// Runs SQL file if nothing exists in the database
func (db *Database) createTables(currentTables []string) {
	if len(currentTables) == 0 {
		// Load DB File
	}
}
