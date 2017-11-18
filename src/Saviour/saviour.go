package main

import (
	"flag"
)

const (
	thisModuleMain = "Main"
)

func main() {
	debugOn := flag.Bool("dbg", false, "Turns On Debug Messages")
	flag.Parse()
	InitDebug(*debugOn)
	DebugHandler.Sys("DebugEnabled", "")
	DebugHandler.Sys("Starting", "")
	InitCron()
	db := InitDatabase()
	InitSystem(db)
}
