/*
This is the beginning entry point for the Saviour Server
*/
package main

import (
  "modules/logger"
  "modules/database"
  "config"
  "fmt"
  "modules/cache"
)

const (
  thisModule = "Main"
)

func main() {
  fmt.Println("Saviour::Starting...")
  fmt.Println("Saviour::LoadingConfiguration")
  conf := config.GetSettings()
  log := logger.InitLogData(conf)
  db := database.InitDatabase(conf, log)
  loadedCache := cache.InitCache(conf, db, log)
  loadedCache.CheckCache()
}