/*
This is the beginning entry point for the Saviour Server
*/
package main

import (
        "modules/logger"
        "modules/database"
        "config"
        "fmt"
)

const (
  thisModule = "Main"
)

func main() {
  fmt.Println("Saviour::Starting...")
  conf := config.GetSettings()
  log := logger.InitLogData(conf)
  db := database.InitDatabase(conf, log)
  db.CheckDB()
}