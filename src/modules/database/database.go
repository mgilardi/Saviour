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
  if err != nil {
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

// WriteCache creates a new cache entry with any object converted into a binary entry
func (db Database) WriteCache(key string, blob []byte, created int64, expires int64) {
  rows, err := db.sql.Query("CALL WriteCache($1, $2, $3, $4)", key, blob, created, expires)
  if err != nil {
    // Error
  }
  rows.Close()
}

// ReadCache finds cache entry and returns the data
func (db Database) ReadCache(key string) (error, []byte) {
  var result sql.NullString
  err := db.sql.QueryRow("CALL ReadCache($1)", key).Scan(&result)
  if err != nil {
    //
  }
  if result.Valid {

  } else {

  }
  convResult := []byte(result.String)
  return err, convResult
}

// RemoveCache finds cache entry and removes it from the database
func (db Database) RemoveCache(key string) {
  rows, err := db.sql.Query("CALL RemoveCache($1)", key)
  if err != nil {
    //
  }
  rows.Close()
}

func (db Database) GetCacheTable() {

}

func (db Database) CreateUser(user string, pass string) {

}

func (db Database) SetUserData(uid int, name string, phone string, address string) {

}

func (db Database) SetPassword(user string) {

}

func (db Database) GetUID(user string) {

}

func (db Database) GetUserData(user string) {

}