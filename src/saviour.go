/*
This is the beginning entry point for the Saviour Server
*/
package main

import (
	"modules/database"
	"modules/logger"
	"modules/system"
)

const (
	thisModule = "Main"
)

func main() {
	logger.SystemMessage("Starting", "")
	db := database.InitDatabase()
	cache := database.InitCache(db)
	system.InitSystem(db, cache)
}
