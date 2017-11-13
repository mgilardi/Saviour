/*
This is the beginning entry point for the Saviour Server
*/
package main

import (
  "modules/database"
  "modules/cache"
  "modules/system"
  "modules/logger"
)

const (
  thisModule = "Main"
)

func main() {
  logger.SystemMessage("Starting", "")
  db := database.InitDatabase()
  loadedCache := cache.InitCache(db)
  system.InitSystem(db, loadedCache)
}