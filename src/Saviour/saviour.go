package main

import (
	"flag"
)

const (
	thisModuleMain = "Main"
)

func main() {
	debugOn := flag.Bool("DebugHandler", false, "Turns On Debug Messages")
	flag.Parse()
	InitDebug(*debugOn)
	DebugHandler.Sys("DebugEnabled", "")
	DebugHandler.Sys("Starting", "")
	db := InitDatabase()
	InitSystem(db)
}
