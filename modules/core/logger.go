// Package core logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written
package core

import (
	"fmt"
	"os"
)

// System Message Constants
const (
	MSG       = "Saviour"
	WARN      = "WARN"
	ERROR     = "ERROR"
	FATAL     = "FATAL"
	MODULELOG = "Logger"
)

var logHandler []log

type log interface {
	SetError(string, string, string)
	CheckLevel() bool
	Enabled() bool
	Write()
}

// Debug contains the enable flag
type debug struct {
	isOn         bool
	debugLevel   int
	currentError string
	currentLevel int
}

// Logger struct contains database
type logger struct {
	logType      string
	logModule    string
	logMsg       string
	currentLevel int
	logLevel     int
	isOn         bool
}

// InitDebug initializes the debug struct
func InitDebug(on bool) {
	var newDebug debug
	newDebug.isOn = on
	logHandler = append(logHandler, &newDebug)
}

// @TODO Remove case statement
func (dbg *debug) SetError(message string, module string, typ string) {
	switch typ {
	case MSG:
		dbg.currentError = typ + "::" + module + "::" + message
		dbg.currentLevel = 3
		handleLog(dbg)
	case WARN:
		dbg.currentError = typ + "::" + module + "::" + message
		dbg.currentLevel = 2
		handleLog(dbg)
	case ERROR:
		dbg.currentError = typ + "::" + module + "::" + message
		dbg.currentLevel = 1
		handleLog(dbg)
	case FATAL:
		dbg.isOn = true
		dbg.currentError = typ + "::" + module + "::" + message
		dbg.currentLevel = 0
		handleLog(dbg)
	default:
		Logger("UnknownErrorType", PACKAGE+"."+MODULELOG+".SetErrorDebug", ERROR)
	}
}

// CheckLevel checks that the current message is above the current log level
func (dbg *debug) CheckLevel() bool {
	// Allow to change log level
	return true
}

// Enabled checks that debug is enabled
func (dbg *debug) Enabled() bool {
	return dbg.isOn
}

// Write writes current system message to the console
func (dbg *debug) Write() {
	fmt.Println(dbg.currentError)
	if dbg.currentLevel == 0 {
		os.Exit(1)
	}
}

// InitLogger constructs logger type
func InitLogger() {
	var newLogDB logger
	newLogDB.isOn = true
	newLogDB.logLevel = 2
	logHandler = append(logHandler, &newLogDB)
}

func (logger *logger) SetError(message string, module string, typ string) {
	switch typ {
	case MSG:
		logger.logType = typ
		logger.logModule = module
		logger.logMsg = message
		logger.currentLevel = 3
		handleLog(logger)
	case WARN:
		logger.logType = typ
		logger.logModule = module
		logger.logMsg = message
		logger.currentLevel = 2
		handleLog(logger)
	case ERROR:
		logger.logType = typ
		logger.logModule = module
		logger.logMsg = message
		logger.currentLevel = 1
		handleLog(logger)
	case FATAL:
		logger.logType = typ
		logger.logModule = module
		logger.logMsg = message
		logger.currentLevel = 0
		handleLog(logger)
	default:
		Logger("UnknownErrorType", PACKAGE+"."+MODULELOG+".SetErrorLogger", ERROR)
	}
}

// CheckLevel checks to make sure current message is above the log level
func (logger *logger) CheckLevel() bool {
	if logger.logLevel > logger.currentLevel {
		return true
	}
	return false
}

// Enabled checks to make sure that the logger module is enabled
func (logger *logger) Enabled() bool {
	return logger.isOn
}

// WriteLog writes log entry into the database
func (logger *logger) Write() {
	writeLog := DBHandler.SetupExec(
		`INSERT INTO logger (type, module, message) `+
			`VALUES (?, ?, ?)`, logger.logType, logger.logModule, logger.logMsg)
	DBHandler.Exec(writeLog)
}

// Logger Global input variable for logger module
func Logger(msg string, module string, typ string) {
	for _, logData := range logHandler {
		logData.SetError(msg, module, typ)
	}
}

// handleLog does checks for both Debug/Logger and writes the output
func handleLog(output log) {
	if output.Enabled() && output.CheckLevel() {
		output.Write()
	}
}
