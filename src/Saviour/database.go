// Package database checks database exists initializes the Database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package main

import (
	"database/sql"
	"errors"
	"time"
	// MySql
	_ "github.com/go-sql-driver/mysql"
)

const (
	thisModuleDB = "Database"
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
	DebugHandler.Sys("Starting", thisModuleDB)
	db.options = GetOptions(thisModuleDB)
	if err != nil {
		DebugHandler.Err(err, thisModuleDB, 1)
	}
	if db.options["User"] == nil {
		DebugHandler.Err(err, thisModuleDB, 1)
	}
	if db.options["Pass"] == nil {
		DebugHandler.Err(err, thisModuleDB, 1)
	}
	user = db.options["User"].(string)
	pass = db.options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	DebugHandler.Sys("DSNLoaded", thisModuleDB)
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		DebugHandler.Err(err, thisModuleDB, 1)
	}
	err = db.sql.Ping()
	if err != nil {
		DebugHandler.Err(err, thisModuleDB, 1)
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
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 1)

	}
	DebugHandler.Sys("AvaliableTables", thisModuleDB)
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			LogHandler.Err(err, thisModuleDB)
			DebugHandler.Err(err, thisModuleDB, 1)
		}
		DebugHandler.Sys(table, thisModuleDB)
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

// CheckUserLogin inputs username and password and checks if it matches
// the database returns a boolean if the user if found along with the
// users database uid
func (db *Database) CheckUserLogin(name string, pass string) (bool, int) {
	var dbPass string
	var uid int
	var verified = false
	DebugHandler.Sys("CheckUserLogin::"+name, thisModuleDB)
	err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	case dbPass == pass:
		verified = true
	default:
		DebugHandler.Sys("InvalidPassword::"+name, thisModuleDB)
		LogHandler.Warn(errors.New("InvalidPassword::"+name), thisModuleDB)
	}
	return verified, uid
}

// CheckUserExits checks the database for a username and returns true or false
// if it exists
func (db *Database) CheckUserExits(name string) bool {
	DebugHandler.Sys("CheckUserExists::"+name, thisModuleDB)
	exists := false
	_, err := db.GetUserID(name)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		exists = true
	}
	return exists
}

// CreateUser creates a new user entry in the database
func (db *Database) CreateUser(name string, pass string, email string) {
	DebugHandler.Sys("CreateUser::"+name, thisModuleDB)
	_, err := db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// RemoveUser removes a user entry from the database
func (db *Database) RemoveUser(name string) {
	DebugHandler.Sys("RemoveUser::"+name, thisModuleDB)
	uid, err := db.GetUserID(name)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	} else {
		_, err := db.sql.Exec(`DELETE FROM * WHERE uid = ?`, uid)
		if err != nil {
			LogHandler.Err(err, thisModuleDB)
			DebugHandler.Err(err, thisModuleDB, 3)
		}
	}
}

// GetUserID will return the database uid for a username
func (db *Database) GetUserID(name string) (int, error) {
	var uid int
	DebugHandler.Sys("GetUserID::"+name, thisModuleDB)
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
	return uid, err
}

// GetUserMap loads user information from the database and returns it
// inside of a map
func (db *Database) GetUserMap(uid int) (map[string]interface{}, error) {
	var userData map[string]interface{}
	var name, email, token string
	err := db.sql.QueryRow(`SELECT name, mail, token FROM users JOIN login_token ON users.uid = login_token.uid AND users.uid = ?`, uid).Scan(&name, &email, &token)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		userData = make(map[string]interface{})
		userData["Name"] = name
		userData["Email"] = email
		userData["Token"] = token
		DebugHandler.Sys("GetUserMap::"+userData["Name"].(string), thisModuleDB)
	}
	return userData, err
}

// CheckToken checks if the user has a token in the database login_token table
func (db *Database) CheckToken(uid int) bool {
	var token string
	exists := false
	DebugHandler.Sys("CheckingToken", thisModuleDB)
	err := db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("TokenNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		DebugHandler.Sys("TokenFound", thisModuleDB)
		exists = true
	}
	return exists
}

// StoreToken writes user token to the database
func (db *Database) StoreToken(uid int, token string) {
	DebugHandler.Sys("StoreToken", thisModuleDB)
	_, err := db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`, uid, token)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// WriteCache creates a new cache entry
func (db *Database) WriteCache(cid string, data []byte) {
	//debug.DebugHandler.Sys("WriteCache", thisModule)
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, data)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// WriteCacheExp creates a new cache entry that expires
func (db *Database) WriteCacheExp(cid string, data []byte, expires int64) {
	//debug.DebugHandler.Sys("WriteCacheExp", thisModule)
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data, expires) VALUES (?, ?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?, expires = ?`, cid, data, expires, data, expires)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// ReadCache returns a cache entry
func (db *Database) ReadCache(cid string) (bool, []byte) {
	var data []byte
	var expires sql.NullInt64
	var exists bool
	DebugHandler.Sys("ReadCache", thisModuleDB)
	err := db.sql.QueryRow(`SELECT data, expires FROM cache WHERE cid = ?`, cid).Scan(&data, &expires)
	exp, _ := expires.Value()
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("CacheNotFound", thisModuleDB)
		exists = false
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
		exists = false
	case expires.Valid && exp.(int64) > time.Now().Unix():
		exists = false
	default:
		exists = true
	}
	return exists, data
}

// ClearCache clears all records in cache table
func (db *Database) ClearCache() {
	_, err := db.sql.Exec(`DELETE FROM cache`)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// CheckCache iterates through cache records checks if the record is expired
// if the record is expired the record is deleted, if the record has NULL for
// a value it is permanent and is skipped, if the record has not expired it is skipped
func (db *Database) CheckCache() {
	var cid string
	var expires sql.NullInt64
	rows, err := db.sql.Query(`SELECT cid, expires FROM cache`)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("ExpiredRecordsNotFound", thisModuleDB)
	case err != nil && err.Error() != "EOF":
		DebugHandler.Err(err, thisModuleDB, 3)
		LogHandler.Err(err, thisModuleDB)
	default:
		for rows.Next() {
			rows.Scan(&cid, &expires)
			if expires.Valid && expires.Int64 < time.Now().Unix() {
				DebugHandler.Sys("RemovingExpired::"+cid, thisModuleDB)
				_, err := db.sql.Exec(`DELETE FROM cache WHERE cid = ?`, cid)
				if err != nil {
					DebugHandler.Err(err, thisModuleDB, 3)
					LogHandler.Err(err, thisModuleDB)
				}
			}
		}
	}
}

// WriteLog writes log entry into the database
func (db *Database) WriteLog(logType string, module string, message string) {
	DebugHandler.Sys("WritingLog::"+logType+"::"+module+"::"+message, thisModuleDB)
	_, err := db.sql.Exec(`INSERT INTO logger (type, module, message) VALUES (?, ?, ?)`, logType, module, message)
	if err != nil {
		DebugHandler.Err(err, thisModuleDB, 3)
		LogHandler.Err(err, thisModuleDB)
	}
}
