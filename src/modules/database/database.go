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
func (db *Database) CheckDB() {
  tables := make([]string, 0)
  rows, err := db.sql.Query(`SHOW TABLES`)
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
func (db *Database) createTables(currentTables []string) {
  if len(currentTables) == 0 {
    // Load DB File
  }
  db.logger.SystemMessage("Tables::Loaded", thisModule)
}

func (db *Database) CheckUserLogin(name string, pass string) (bool,int) {
  var dbPass string
  var uid int
  var verified = false
  err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`,name).Scan(&dbPass, &uid)
  switch {
  case err == sql.ErrNoRows:
    db.logger.Error("UserNotFound", thisModule, 3)
  case err != nil:
    db.logger.Error(err.Error(), thisModule, 3)
  case dbPass == pass:
    verified = true
  }
  return verified, uid
}

func (db *Database) CheckUserExits(name string) bool {
  exists := false
  _ , err := db.GetUserID(name)
  switch {
  case err == sql.ErrNoRows:
    db.logger.Error("UserNotFound", thisModule, 3)
  case err != nil:
    db.logger.Error(err.Error(), thisModule, 3)
  default:
    exists = true
  }
  return exists
}

func (db *Database) CreateUser(name string, pass string, email string) {
  _, err := db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
    if (err != nil) {
      db.logger.Error(err.Error(), thisModule, 3)
    }
}

func (db *Database) RemoveUser(name string) {
  uid , err := db.GetUserID(name)
  rows, err := db.sql.Query(`DELETE FROM * WHERE uid = ?`, uid)
  if err != nil {
    db.logger.Error(err.Error(), thisModule, 3)
  }
  rows.Close()
}

func (db *Database) GetUserID(name string) (int, error) {
  var uid int
  err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
  if (err !=nil) {
    db.logger.Error(err.Error(), thisModule, 3)
  }
  return uid , err
}

func (db *Database) GetUserMap(uid int) (map[string]string, error) {
  var userData map[string]string
  var name, email, token string
  err := db.sql.QueryRow(`SELECT name, mail, token FROM users JOIN login_token ON users.uid = login_token.uid AND users.uid = ?`, uid).Scan(&name, &email, &token)
  if err != nil {
    db.logger.Error("GetUserMap::" + err.Error(), thisModule, 3)
  }
  userData = make(map[string]string)
  userData["name"] = name
  userData["email"] = email
  userData["token"] = token
  return userData, err
}

func (db *Database) CheckToken(uid int) bool {
  exists := false
  var token string
  err := db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
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

func (db *Database) StoreToken(uid int, token string) {
  _, err := db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`, uid, token)
  if (err != nil) {
    db.logger.Error(err.Error(), thisModule, 3)
  }
}

func WriteCache() {

}

func ReadCache() {

}