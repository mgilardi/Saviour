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
	core.Logger("DebugEnabled", "", core.MSG)
	core.Logger("Starting", "", core.MSG)
	core.InitOptions()
	core.InitCron()
	core.InitDatabase()
	core.InitLogger()
	core.InitCache()
	core.InitAccess()
	core.InitSystem()
}
