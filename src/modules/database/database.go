// Package database checks database exists initializes the Database object provides
// object methods for reading/writing/managing of the database
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

type Database struct {
  sql *sql.DB
  options *config.Setting
  logger *logger.LogData
}

// InitDatabase initialize the database object and passes a pointer to the main loop
func InitDatabase(settings *[]config.Setting, log *logger.LogData) *Database {
  var db Database
  var err error
  var user, pass, dsn string
  db.logger = log
  err, db.options = config.GetSettingModule(thisModule, settings)
  err, user = db.options.FindValue("User")
  err, pass = db.options.FindValue("Pass")
  if (err != nil) {
    db.logger.Error("CannotLoadSettings", thisModule, 1)
    db.logger.Error(err.Error(), thisModule, 3)
  }
  dsn = user + ":" + pass + "@/"
  // Open Database
  db.sql, err = sql.Open("mysql", dsn)
  if (err != nil) {
    db.logger.Error("CannotOpenDB", thisModule, 1)
    db.logger.Error(err.Error(), thisModule, 3)
  }
  err = db.sql.Ping()
  if err != nil {
    db.logger.Error("CannotPingDB", thisModule, 1)
  }
  return &db
}

// CheckDB checks if the saviour database and tables exist if not it creates them.
func (db Database) CheckDB() {
  var name string
  var exists bool
  rows, err := db.sql.Query("show databases")
  if err != nil {
    db.logger.Error("CannotQueryDB", thisModule, 1)
    db.logger.Error(err.Error(), thisModule, 3)
  }
  exists = false
  for rows.Next() {
    rows.Scan(&name)
    if (name == "saviour") {
      exists = true
      db.logger.SystemMessage("DBFound", thisModule, 1)
      break
    }
  }
  if !exists {
    db.logger.Error("DBDoesNotExist", thisModule, 2)
  }
}
