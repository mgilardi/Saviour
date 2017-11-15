// Package database checks database exists initializes the Database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package database

import (
	"config"
	"database/sql"
	"errors"
	"modules/debug"
	"time"
	// MySql
	_ "github.com/go-sql-driver/mysql"
)

const (
	thisModule = "Database"
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
	debug.Dbg.Sys("Starting", thisModule)
	db.options = config.GetOptions(thisModule)
	if err != nil {
		debug.Dbg.Err(err, thisModule, 1)
	}
	if db.options["User"] == nil {
		debug.Dbg.Err(err, thisModule, 1)
	}
	if db.options["Pass"] == nil {
		debug.Dbg.Err(err, thisModule, 1)
	}
	user = db.options["User"].(string)
	pass = db.options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	debug.Dbg.Sys("DSNLoaded", thisModule)
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		debug.Dbg.Err(err, thisModule, 1)
	}
	err = db.sql.Ping()
	if err != nil {
		debug.Dbg.Err(err, thisModule, 1)
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
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 1)

	}
	debug.Dbg.Sys("AvaliableTables", thisModule)
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			debug.Dbg.Err(err, thisModule, 1)
			LogDB.Err(err, thisModule)
		}
		debug.Dbg.Sys(table, thisModule)
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
	debug.Dbg.Sys("CheckUserLogin::"+name, thisModule)
	err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		debug.Dbg.Sys("UserNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	case dbPass == pass:
		verified = true
	default:
		debug.Dbg.Sys("InvalidPassword::"+name, thisModule)
		LogDB.Warn(errors.New("InvalidPassword::"+name), thisModule)
	}
	return verified, uid
}

// CheckUserExits checks the database for a username and returns true or false
// if it exists
func (db *Database) CheckUserExits(name string) bool {
	debug.Dbg.Sys("CheckUserExists::"+name, thisModule)
	exists := false
	_, err := db.GetUserID(name)
	switch {
	case err == sql.ErrNoRows:
		debug.Dbg.Sys("UserNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	default:
		exists = true
	}
	return exists
}

// CreateUser creates a new user entry in the database
func (db *Database) CreateUser(name string, pass string, email string) {
	debug.Dbg.Sys("CreateUser::"+name, thisModule)
	_, err := db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
}

// RemoveUser removes a user entry from the database
func (db *Database) RemoveUser(name string) {
	debug.Dbg.Sys("RemoveUser::"+name, thisModule)
	uid, err := db.GetUserID(name)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	} else {
		_, err := db.sql.Exec(`DELETE FROM * WHERE uid = ?`, uid)
		if err != nil {
			LogDB.Err(err, thisModule)
			debug.Dbg.Err(err, thisModule, 3)
		}
	}
}

// GetUserID will return the database uid for a username
func (db *Database) GetUserID(name string) (int, error) {
	var uid int
	debug.Dbg.Sys("GetUserID::"+name, thisModule)
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		debug.Dbg.Sys("UserNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
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
		debug.Dbg.Sys("UserNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	default:
		userData = make(map[string]interface{})
		userData["Name"] = name
		userData["Email"] = email
		userData["Token"] = token
		debug.Dbg.Sys("GetUserMap::"+userData["Name"].(string), thisModule)
	}
	return userData, err
}

// CheckToken checks if the user has a token in the database login_token table
func (db *Database) CheckToken(uid int) bool {
	var token string
	exists := false
	debug.Dbg.Sys("CheckingToken", thisModule)
	err := db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		debug.Dbg.Sys("TokenNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	default:
		debug.Dbg.Sys("TokenFound", thisModule)
		exists = true
	}
	return exists
}

// StoreToken writes user token to the database
func (db *Database) StoreToken(uid int, token string) {
	debug.Dbg.Sys("StoreToken", thisModule)
	_, err := db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`, uid, token)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
}

// WriteCache creates a new cache entry
func (db *Database) WriteCache(cid string, data []byte) {
	//debug.Dbg.Sys("WriteCache", thisModule)
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, data)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
}

// WriteCacheExp creates a new cache entry that expires
func (db *Database) WriteCacheExp(cid string, data []byte, expires int64) {
	//debug.Dbg.Sys("WriteCacheExp", thisModule)
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data, expires) VALUES (?, ?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, expires, data)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
}

// ReadCache returns a cache entry
func (db *Database) ReadCache(cid string) []byte {
	var data []byte
	debug.Dbg.Sys("ReadCache", thisModule)
	err := db.sql.QueryRow(`SELECT data FROM cache WHERE cid = ?`, cid).Scan(&data)
	switch {
	case err == sql.ErrNoRows:
		debug.Dbg.Sys("CacheNotFound", thisModule)
	case err != nil:
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
	return data
}

// ClearCache clears all records in cache table
func (db *Database) ClearCache() {
	_, err := db.sql.Exec(`DELETE FROM cache`)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
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
		debug.Dbg.Sys("ExpiredRecordsNotFound", thisModule)
	case err != nil && err.Error() != "EOF":
		debug.Dbg.Err(err, thisModule, 3)
		LogDB.Err(err, thisModule)
	default:
		for rows.Next() {
			rows.Scan(&cid, &expires)
			if expires.Valid && expires.Int64 < time.Now().Unix() {
				debug.Dbg.Sys("RemovingExpired::"+cid, thisModule)
				_, err := db.sql.Exec(`DELETE FROM cache WHERE cid = ?`, cid)
				if err != nil {
					debug.Dbg.Err(err, thisModule, 3)
					LogDB.Err(err, thisModule)
				}
			}
		}
	}
}

// WriteLog writes log entry into the database
func (db *Database) WriteLog(logType string, module string, message string) {
	debug.Dbg.Sys("WritingLog::"+logType+"::"+module+"::"+message, thisModule)
	_, err := db.sql.Exec(`INSERT INTO logger (type, module, message) VALUES (?, ?, ?)`, logType, module, message)
	if err != nil {
		debug.Dbg.Err(err, thisModule, 3)
		LogDB.Err(err, thisModule)
	}
}
