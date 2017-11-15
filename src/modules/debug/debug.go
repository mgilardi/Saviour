// Package logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written
package debug

import (
	"fmt"
	"log"
)

const (
	thisModule      = "debug"
	typeRequest     = "LogType"
	logLevelRequest = "LogLevel"
	sysLevelRequest = "SysLevel"
)

// Dbg global variable for debug access
var Dbg *Debug

// Debug contains the enable flag
type Debug struct {
	isOn bool
}

// InitDebug initializes the debug struct
func InitDebug(on bool) {
	var debug Debug
	debug.isOn = on
	Dbg = &debug
}

// Sys outputs system messsages in the console when dbg is enabled
func (dbg *Debug) Sys(message string, module string) {
	if dbg.isOn {
		fmt.Println("Saviour::" + module + "::" + message)
	}
}

// Err outputs a assembled error message to the console when dbg is enabled
// is the error level is 1 then a fatal error is called which forces the
// application closed
func (dbg *Debug) Err(err error, module string, level int) {
	switch {
	case dbg.isOn:
		fmt.Println("Error::" + module + "::" + err.Error())
	case level == 1:
		log.Fatal(err)
	default:
		// Ignore
	}
}
