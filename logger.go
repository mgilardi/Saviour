package main

// logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written

const (
	logLevelRequest = "LogLevel"
	sysLevelRequest = "SysLevel"
)

// LogHandler logger global variable
var LogHandler *Logger

// Logger struct contains database
type Logger struct {
	db *Database
}

// InitLogger constructs logger type
func InitLogger(loadDB *Database) {
	var newLogDB Logger
	newLogDB.db = loadDB
	LogHandler = &newLogDB
}

// Stat writes status message
func (logger *Logger) Stat(message string, module string) {
	logger.db.WriteLog("Status", module, message)
}

// Warn writes warning messages
func (logger *Logger) Warn(err error, module string) {
	logger.db.WriteLog("Warn", module, err.Error())
}

// Err writes error messages
func (logger *Logger) Err(err error, module string) {
	logger.db.WriteLog("Error", module, err.Error())
}

// Fatal Writes Fatal Messages
func (logger *Logger) Fatal(err error, module string) {
	logger.db.WriteLog("Fatal", module, err.Error())
}
