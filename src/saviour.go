/*
This is the beginning entry point for the Saviour Server
*/
package main

import (
	"flag"
	"modules/database"
	"modules/debug"
	"modules/system"
)

const (
	thisModule = "Main"
)

func main() {
	debugOn := flag.Bool("dbg", false, "Turns On Debug Messages")
	flag.Parse()
	debug.InitDebug(*debugOn)
	debug.Dbg.Sys("DebugEnabled", "")
	debug.Dbg.Sys("Starting", "")
	db := database.InitDatabase()
	system.InitSystem(db)
}
