package main

// @TODO Check for existing DB and if it doesn't exist create it and store password in settings file.
// @TODO Explore debugging/break points
// @TODO Set width standard for coding + find readme.md file.

import (
	"Saviour/modules/core"
	"flag"
)

const (
	thisModuleMain = "Main"
)

func main() {
	debugOn := flag.Bool("dbg", false, "Turns On Debug Messages")
	inProd := flag.Bool("prod", false, "Turns On Production Mode")
	flag.Parse()
	core.InitDebug(*debugOn)
	core.Logger("DebugEnabled", "", core.MSG)
	core.Logger("Starting", "", core.MSG)
	core.InitOptions()
	core.InitCron()
	core.InitDatabase()
	core.InitLogger()
	core.InitCache()
	core.InitAgent()
	core.InitAccess()
	core.InitSystem(*inProd)
	core.InitCommand()
}
