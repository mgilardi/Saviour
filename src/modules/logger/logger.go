package logger

import (
  "fmt"
  "config"
  "strconv"
)

const (
  thisModule = "Logger"
  typeRequest = "FileType"
  logLevelRequest = "LogLevel"
  sysLevelRequest = "SysLevel"
)

type LogData struct {
  logMessageLevel int
  systemMessageLevel int
  logType string
  options *config.Setting
}

// InitLogData 
func InitLogData(settings *[]config.Setting) *LogData {
  var log LogData
  var buf string
  err, options := config.GetSettingModule(thisModule, settings)
  if (err != nil) {
    log.Error("CouldNotLoadSettings", thisModule, 1)
  }
  log.options = options
  err, log.logType = log.options.FindValue(typeRequest)
  err, buf = log.options.FindValue(sysLevelRequest)
  log.systemMessageLevel, err = strconv.Atoi(buf)
  err, buf = log.options.FindValue(logLevelRequest)
  log.logMessageLevel, err = strconv.Atoi(buf)
  if (err != nil) {
    log.Error("CannotLoadConfig:" + sysLevelRequest, thisModule, 1)
    log.Error(err.Error(), thisModule, 3)
  }
  return &log
}

func (log *LogData) SystemMessage(message string, module string, level int) {
  if (log.systemMessageLevel < level) {
    // ignore message
  } else {
    fmt.Println("Saviour::" + module + "::" + message)
  }
}

func (log *LogData) Error(message string, module string, level int) {
  if (log.logMessageLevel < level) {
    // ignore message
  } else {
    fmt.Println("Error::" + module + "::" + message)
  }
}