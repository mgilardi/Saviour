// Package core logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written
package core

import (
	"fmt"
	"os"
)

type log interface {
	SetError(string, string, int)
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

// System Message Constants
const (
	MSG   = 0
	WARN  = 1
	ERROR = 2
	FATAL = 3
)

var logHandler []log

// InitDebug initializes the debug struct
func InitDebug(on bool) {
	var newDebug debug
	newDebug.isOn = on
	logHandler = append(logHandler, &newDebug)
}

// @TODO Remove case statement
func (dbg *debug) SetError(message string, module string, typ int) {
	switch typ {
	case 0:
		dbg.Sys(message, module)
	case 1:
		dbg.Warn(message, module)
	case 2:
		dbg.Err(message, module)
	case 3:
		dbg.Fatal(message, module)
		os.Exit(1)
	}
}

// Sys outputs system messsages in the console when dbg is enabled
func (dbg *debug) Sys(message string, module string) {
	dbg.currentError = "Saviour::" + module + "::" + message
	dbg.currentLevel = 3
	handleLog(dbg)
}

// Warn writes warning messages
func (dbg *debug) Warn(err string, module string) {
	dbg.currentError = "Warn::" + module + "::" + err
	dbg.currentLevel = 2
	handleLog(dbg)
}

// Err writes error messages
func (dbg *debug) Err(err string, module string) {
	dbg.currentError = "Error::" + module + "::" + err
	dbg.currentLevel = 1
	handleLog(dbg)
}

// Fatal Writes Fatal Messages
func (dbg *debug) Fatal(err string, module string) {
	dbg.currentError = "Fatal::" + module + "::" + err
	dbg.currentLevel = 0
	dbg.isOn = true
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
	logType      string
	logModule    string
	logMsg       string
	currentLevel int
	logLevel     int
	isOn         bool
}

// InitLogger constructs logger type
func InitLogger() {
	var newLogDB logger
	newLogDB.isOn = true
	newLogDB.logLevel = 2
	logHandler = append(logHandler, &newLogDB)
}

func (logger *logger) SetError(message string, module string, typ int) {
	switch typ {
	case 0:
		logger.Sys(message, module)
	case 1:
		logger.Warn(message, module)
	case 2:
		logger.Err(message, module)
	case 3:
		logger.Fatal(message, module)
	}
}

// Sys writes status message
func (logger *logger) Sys(message string, module string) {
	logger.logType = "Status"
	logger.logModule = module
	logger.logMsg = message
	logger.currentLevel = 3
	handleLog(logger)
}

// Warn writes warning messages
func (logger *logger) Warn(err string, module string) {
	logger.logType = "Warn"
	logger.logModule = module
	logger.logMsg = err
	logger.currentLevel = 2
	handleLog(logger)
}

// Err writes error messages
func (logger *logger) Err(err string, module string) {
	logger.logType = "Error"
	logger.logModule = module
	logger.logMsg = err
	logger.currentLevel = 1
	handleLog(logger)
}

// Fatal Writes Fatal Messages
func (logger *logger) Fatal(err string, module string) {
	logger.logType = "Fatal"
	logger.logModule = module
	logger.logMsg = err
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
	_, err := DBHandler.sql.Exec(
		`INSERT INTO logger (type, module, message) VALUES (?, ?, ?)`,
		logger.logType, logger.logModule, logger.logMsg)
	if err != nil {
		Logger(err.Error(), "Logger", ERROR)
	}
}

// Logger Global input variable for logger module
func Logger(msg string, module string, typ int) {
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
