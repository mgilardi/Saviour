package main

import (
	"Saviour/modules/core"
	"flag"
)

const (
	thisModuleMain = "Main"
)

func main() {
	debugOn := flag.Bool("dbg", false, "Turns On Debug Messages")
	flag.Parse()
	core.InitDebug(*debugOn)
	core.DebugHandler.Sys("DebugEnabled", "")
	core.DebugHandler.Sys("Starting", "")
	core.InitCron()
	db := core.InitDatabase()
	core.InitSystem(db)
}
