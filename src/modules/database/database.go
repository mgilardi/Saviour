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

func (db Database) createTables(currentTables []string) {
  if len(currentTables) == 0 {
    rows, err := db.sql.Query("CREATE TABLE Cache (CID INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
      " Data LONGBLOB, Created INT(11), Expires INT(11))")
    if err != nil {
        db.logger.Error(err.Error(), thisModule, 3)
        db.logger.Error("FailedToLoadCacheTable", thisModule, 1)
    }
    rows, err = db.sql.Query("CREATE TABLE Users (UID INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
      " Username VARCHAR(45), Password VARCHAR(45), FirstName VARCHAR(45), LastName VARCHAR(45)," +
      " Phone VARCHAR(45), Email VARCHAR(45), Address VARCHAR(45))")
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("FailedToLoadUsersTable", thisModule, 1)
    }
    rows, err = db.sql.Query("CREATE TABLE CreditCards (CCID INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
      " UserID INT NOT NULL, Type VARCHAR(45), Number INT(45)," +
      " CCV INT(45), INDEX CreditCard_User_IDX(UserID))")
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("FailedToLoadCreditCardsTable", thisModule, 1)
    }
    rows, err = db.sql.Query("ALTER TABLE `CreditCards` ADD CONSTRAINT `userID_creditcards` FOREIGN KEY (`UserID`)" +
      "REFERENCES `Users`(`UID`) ON DELETE CASCADE ON UPDATE CASCADE")
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("FailedToLoadCreditCards", thisModule, 1)
      }
    rows, err = db.sql.Query("CREATE TABLE ActivityLog (AID INT NOT NULL AUTO_INCREMENT PRIMARY KEY," +
      " UserID INT NOT NULL, Timestamp DATETIME," +
      " Log VARCHAR(255), INDEX CreditCard_User_IDX(UserID))")
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("FailedToLoadActivityLogTable", thisModule, 1)
    }
    rows, err = db.sql.Query("ALTER TABLE `ActivityLog` ADD CONSTRAINT `userID_activityLog` FOREIGN KEY (`UserID`)" +
      "REFERENCES `Users`(`UID`) ON DELETE CASCADE ON UPDATE CASCADE")
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("FailedToLoadActivityLogTable", thisModule, 1)
    }
    err = rows.Err()
    if err != nil {
      db.logger.Error(err.Error(), thisModule, 3)
      db.logger.Error("RowsError", thisModule, 1)
    }
  }
  db.logger.SystemMessage("Tables::Loaded", thisModule)
}