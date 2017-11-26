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
	core.Sys("DebugEnabled", "")
	core.Sys("Starting", "")
	core.InitOptions()
	core.InitCron()
	core.InitDatabase()
	core.InitLogger()
	core.InitCache()
	core.InitSystem()
}
