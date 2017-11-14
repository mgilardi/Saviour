// Package database checks database exists initializes the Database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package database

import (
	"config"
	"database/sql"
	"modules/logger"
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
	logger.SystemMessage("Starting", thisModule)
	db.options = config.GetOptions(thisModule)
	if err != nil {
		logger.Error("CannotRetrieveSettingsModules", thisModule, 1)
		logger.Error(err.Error(), thisModule, 3)
	}
	if db.options["User"] == nil {
		logger.Error("UsernameNotFound", thisModule, 1)
	}
	if db.options["Pass"] == nil {
		logger.Error("PasswordNotFound", thisModule, 1)
	}
	user = db.options["User"].(string)
	pass = db.options["Pass"].(string)
	db.dsn = user + ":" + pass + "@/saviour"
	logger.SystemMessage("DSNLoaded", thisModule)
	// Open Database
	db.sql, err = sql.Open("mysql", db.dsn)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
		logger.Error("CannotOpenDB", thisModule, 1)
	}
	err = db.sql.Ping()
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
		logger.Error("CannotPingDB", thisModule, 1)
	}
	db.CheckDB()
	return &db
}

// CheckDB checks if database exists and outputs tables that are found.
func (db *Database) CheckDB() {
	tables := make([]string, 0)
	rows, err := db.sql.Query(`SHOW TABLES`)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
		logger.Error("NoTablesCheckDB", thisModule, 1)
	}
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			logger.Error(err.Error(), thisModule, 3)
			logger.Error("CouldNotCheckTables", thisModule, 1)
		}
		logger.SystemMessage("LoadingTable::"+table, thisModule)
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
	logger.SystemMessage("Tables::Loaded", thisModule)
}

// CheckUserLogin inputs username and password and checks if it matches
// the database returns a boolean if the user if found along with the
// users database uid
func (db *Database) CheckUserLogin(name string, pass string) (bool, int) {
	var dbPass string
	var uid int
	var verified = false
	err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		logger.Error("UserNotFound", thisModule, 3)
	case err != nil:
		logger.Error(err.Error(), thisModule, 3)
	case dbPass == pass:
		verified = true
	}
	return verified, uid
}

// CheckUserExits checks the database for a username and returns true or false
// if it exists
func (db *Database) CheckUserExits(name string) bool {
	exists := false
	_, err := db.GetUserID(name)
	switch {
	case err == sql.ErrNoRows:
		logger.Error("UserNotFound", thisModule, 3)
	case err != nil:
		logger.Error(err.Error(), thisModule, 3)
	default:
		exists = true
	}
	return exists
}

// CreateUser creates a new user entry in the database
func (db *Database) CreateUser(name string, pass string, email string) {
	_, err := db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
}

// RemoveUser removes a user entry from the database
func (db *Database) RemoveUser(name string) {
	uid, err := db.GetUserID(name)
	rows, err := db.sql.Query(`DELETE FROM * WHERE uid = ?`, uid)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
	rows.Close()
}

// GetUserID will return the database uid for a username
func (db *Database) GetUserID(name string) (int, error) {
	var uid int
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
	return uid, err
}

// GetUserMap loads user information from the database and returns it
// inside of a map
func (db *Database) GetUserMap(uid int) (map[string]interface{}, error) {
	var userData map[string]interface{}
	var name, email, token string
	err := db.sql.QueryRow(`SELECT name, mail, token FROM users JOIN login_token ON users.uid = login_token.uid AND users.uid = ?`, uid).Scan(&name, &email, &token)
	if err != nil {
		logger.Error("GetUserMap::"+err.Error(), thisModule, 3)
	}
	userData = make(map[string]interface{})
	userData["Name"] = name
	userData["Email"] = email
	userData["Token"] = token
	return userData, err
}

// CheckToken checks if the user has a token in the database login_token table
func (db *Database) CheckToken(uid int) bool {
	exists := false
	var token string
	err := db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		logger.Error("TokenNotFound", thisModule, 3)
	case err != nil:
		logger.Error(err.Error(), thisModule, 3)
	default:
		exists = true
	}
	return exists
}

// StoreToken writes user token to the database
func (db *Database) StoreToken(uid int, token string) {
	_, err := db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`, uid, token)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
}

// WriteCache creates a new cache entry
func (db *Database) WriteCache(cid string, data []byte) {
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, data)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
}

// WriteCacheExp creates a new cache entry that expires
func (db *Database) WriteCacheExp(cid string, data []byte, expires int64) {
	_, err := db.sql.Exec(`INSERT INTO cache (cid, data, expires) VALUES (?, ?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, expires, data)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
}

// ReadCache returns a cache entry
func (db *Database) ReadCache(cid string) []byte {
	var data []byte
	err := db.sql.QueryRow(`SELECT data FROM cache WHERE cid = ?`, cid).Scan(&data)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
	return data
}

// ClearCache clears all records in cache table
func (db *Database) ClearCache() {
	_, err := db.sql.Exec(`DELETE FROM cache`)
	if err != nil {
		//
	}
}

// CheckCache iterates through cache records checks if the record is expired
// if the record is expired the record is deleted, if the record has NULL for
// a value it is permanent and is skipped, if the record has not expired it is skipped
func (db *Database) CheckCache() {
	var cid string
	var expires sql.NullInt64
	rows, err := db.sql.Query(`SELECT cid, expires FROM cache`)
	if err != nil {
		//
	}
	for rows.Next() {
		rows.Scan(&cid, &expires)
		if expires.Valid && expires.Int64 < time.Now().Unix() {
			_, err := db.sql.Exec(`DELETE FROM cache WHERE cid = ?`, cid)
			if err != nil {
				//
			}
		}
	}
}
