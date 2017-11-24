// Package core logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written
package core

import (
	"fmt"
	"os"
)

var logHandler []log

type log interface {
	Sys(string, string)
	Warn(error, string)
	Err(error, string)
	Fatal(error, string)
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

// InitDebug initializes the debug struct
func InitDebug(on bool) {
	var newDebug debug
	newDebug.isOn = on
	logHandler = append(logHandler, &newDebug)
}

// Sys outputs system messsages in the console when dbg is enabled
func (dbg *debug) Sys(message string, module string) {
	dbg.currentError = "Saviour::" + module + "::" + message
	dbg.currentLevel = 2
	handleLog(dbg)
}

// Warn writes warning messages
func (dbg *debug) Warn(err error, module string) {
	dbg.currentError = "Warn::" + module + "::" + err.Error()
	dbg.currentLevel = 3
	handleLog(dbg)
}

// Err writes error messages
func (dbg *debug) Err(err error, module string) {
	dbg.currentError = "Error::" + module + "::" + err.Error()
	dbg.currentLevel = 1
	handleLog(dbg)
}

// Fatal Writes Fatal Messages
func (dbg *debug) Fatal(err error, module string) {
	dbg.currentError = "Fatal::" + module + "::" + err.Error()
	dbg.currentLevel = 0
	handleLog(dbg)
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
}

// Logger struct contains database
type logger struct {
	db           *Database
	logType      string
	logModule    string
	logMsg       string
	currentLevel int
	logLevel     int
	isOn         bool
}

// InitLogger constructs logger type
func InitLogger(loadDB *Database) {
	var newLogDB logger
	newLogDB.isOn = true
	newLogDB.logLevel = 1
	newLogDB.db = loadDB
	logHandler = append(logHandler, &newLogDB)
}

// Sys writes status message
func (logger *logger) Sys(message string, module string) {
	logger.logType = "Status"
	logger.logModule = module
	logger.logMsg = message
	logger.currentLevel = 2
	handleLog(logger)
}

// Warn writes warning messages
func (logger *logger) Warn(err error, module string) {
	logger.logType = "Warn"
	logger.logModule = module
	logger.logMsg = err.Error()
	logger.currentLevel = 3
	handleLog(logger)
}

// Err writes error messages
func (logger *logger) Err(err error, module string) {
	logger.logType = "Error"
	logger.logModule = module
	logger.logMsg = err.Error()
	logger.currentLevel = 1
	handleLog(logger)
}

// Fatal Writes Fatal Messages
func (logger *logger) Fatal(err error, module string) {
	logger.logType = "Fatal"
	logger.logModule = module
	logger.logMsg = err.Error()
	logger.currentLevel = 0
	handleLog(logger)
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
	Sys("WritingLog::"+logger.logType+"::"+logger.logModule+"::"+logger.logMsg, "Logger")
	_, err := logger.db.sql.Exec(`INSERT INTO logger (type, module, message) VALUES (?, ?, ?)`,
		logger.logType, logger.logModule, logger.logMsg)
	if err != nil {
		Error(err, "Cache")
	}
}

// Error outputs an error to both debug and logger
func Error(err error, module string) {
	for _, out := range logHandler {
		out.Err(err, module)
	}
}

// Sys output a system message to both debug and logger
func Sys(msg string, module string) {
	for _, out := range logHandler {
		out.Sys(msg, module)
	}
}

// Warn outputs a warning message to both debug and logger
func Warn(err error, module string) {
	for _, out := range logHandler {
		out.Warn(err, module)
	}
}

// Fatal outputs a fatal message to both debug and logger
func Fatal(err error, module string) {
	for _, out := range logHandler {
		out.Fatal(err, module)
	}
	os.Exit(1)
}

// handleLog does checks for both Debug/Logger and writes the output
func handleLog(out log) {
	if out.Enabled() && out.CheckLevel() {
		out.Write()
	}
}
