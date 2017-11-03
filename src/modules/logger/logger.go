package logger

import (
  "fmt"
  "config"
  "os"
)

const (
  thisModule = "Logger"
  typeRequest = "LogType"
  logLevelRequest = "LogLevel"
  sysLevelRequest = "SysLevel"
)

// LogData object contains
type LogData struct {
  logMessageLevel int
  systemMessageLevel int
  logType string
  options *config.Setting
}

// InitLogData loads the LogData object and returns a pointer to that object
func InitLogData(settings *[]config.Setting) *LogData {
  var log LogData
  log.SystemMessage("Starting", thisModule)
  err, options := config.GetSettingModule(thisModule, settings)
  if (err != nil) {
    log.Error(err.Error(), thisModule, 3)
    log.Error("CouldNotLoadSettings", thisModule, 1)
  }
  log.options = options
  log.logType = log.options.FindValue(typeRequest).(string)
  log.systemMessageLevel = int(log.options.FindValue(sysLevelRequest).(float64))
  log.logMessageLevel = int(log.options.FindValue(logLevelRequest).(float64))
  return &log
}

// SystemMessage outputs a assembled system message to the console
func (log *LogData) SystemMessage(message string, module string) {
  fmt.Println("Saviour::" + module + "::" + message)
}

// Error outputs a assembled error message to the console
func (log *LogData) Error(message string, module string, level int) {
  fmt.Println("Error::" + module + "::" + message)
  if (level == 1) {
    os.Exit(1)
  }
}