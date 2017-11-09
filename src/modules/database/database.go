// Package database checks database exists initializes the Database object provides
// object methods for reading/writing/managing of the database.
// the database is limited to pre-defined actions in the form of routines for security purposes
package database

import (
        "config"
        "database/sql"
        _ "github.com/go-sql-driver/mysql"
        "modules/logger"
)

const (
  thisModule = "Database"
)

// Database tyoe contains the sql access, options, logger, and the dsn for sql login
type Database struct {
  sql *sql.DB
  options *config.Setting
  logger *logger.LogData
  dsn string
}

// InitDatabase initialize the database object and passes a pointer to the main loop
func InitDatabase(settings *[]config.Setting, log *logger.LogData) *Database {
  var db Database
  var err error
  var user, pass string
  db.logger = log
  db.logger.SystemMessage("Starting", thisModule)
  err, db.options = config.GetSettingModule(thisModule, settings)
  if err != nil {
    db.logger.Error("CannotRetrieveSettingsModules", thisModule, 1)
    db.logger.Error(err.Error(), thisModule, 3)
  }
  if db.options.FindValue("User") == nil {
    log.Error("UsernameNotFound",thisModule, 1)
  }
  if db.options.FindValue("Pass") == nil {
    log.Error("PasswordNotFound",thisModule, 1)
  }
  user = db.options.FindValue("User").(string)
  pass = db.options.FindValue("Pass").(string)
  db.dsn = user + ":" + pass + "@/saviour"
  db.logger.SystemMessage("DSNLoaded", thisModule)
  // Open Database
  db.sql, err = sql.Open("mysql", db.dsn)
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
    db.logger.Error("CannotOpenDB", thisModule, 1)
  }
  err = db.sql.Ping()
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
    db.logger.Error("CannotPingDB", thisModule, 1)
  }
  db.CheckDB()
  return &db
}

// CheckDB checks if database exists and outputs tables that are found.
func (db Database) CheckDB() {
  tables := make([]string, 0)
  rows, err := db.sql.Query("SHOW TABLES")
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
    db.logger.Error("NoTablesCheckDB", thisModule, 1)
  }
  for rows.Next() {
    var table string
    err = rows.Scan(&table)
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("CouldNotCheckTables", thisModule, 1)
    }
    db.logger.SystemMessage("LoadingTable::" + table, thisModule)
    tables = append(tables, table)
  }
  db.createTables(tables)
  rows.Close()
}

// Runs SQL file if nothing exists in the database
func (db Database) createTables(currentTables []string) {
  if len(currentTables) == 0 {
    // Load DB File
  }
  db.logger.SystemMessage("Tables::Loaded", thisModule)
}

func (db Database) CheckUser(name string, pass string) bool {
  var dbPass string
  var verified = false
  err := db.sql.QueryRow("SELECT pass FROM users WHERE name = ?",name).Scan(&dbPass)
  switch {
  case err == sql.ErrNoRows:
    db.logger.Error("UserNotFound", thisModule, 3)
  case err != nil:
    db.logger.Error(err.Error(), thisModule, 3)
  case dbPass == pass:
    verified = true
  }
  return verified
}

func (db Database) CreateUser(name string, pass string, email string) {
  rows, err := db.sql.Query("INSERT INTO users (name, pass, mail, created) VALUES (?, ?, ?, NOW())", name, pass, email)
    if (err != nil) {
      db.logger.Error(err.Error(), thisModule, 3)
    }
  rows.Close()
}

func (db Database) GetUserID(name string) int {
  var uid int
  err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
  if (err !=nil) {
    db.logger.Error(err.Error(), thisModule, 3)
  }
  return uid
}

func (db Database) CheckToken(name string) bool {
  exists := false
  var token string
  uid := db.GetUserID(name)
  err := db.sql.QueryRow("SELECT token FROM login_token where uid = ?", uid).Scan(&token)
  switch {
  case err == sql.ErrNoRows:
    db.logger.Error("TokenNotFound", thisModule, 3)
  case err != nil:
    db.logger.Error(err.Error(), thisModule, 3)
  default:
    exists = true
  }
  return exists
}

func (db Database) StoreToken(name string, token string) {
  uid := db.GetUserID(name)
  rows, err := db.sql.Query("INSERT INTO login_token(uid, token, created) VALUES (?, ?, NOW())", uid, token)
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
  }
  rows.Close()
}

func (db Database) GetToken(name string) string {
  var token string
  uid := db.GetUserID(name)
  err := db.sql.QueryRow("SELECT token FROM login_token WHERE uid = ?", uid).Scan(&token)
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
  }
  return token
}

func WriteCache() {

}

func ReadCache() {

}