// Package logger handles all system output, both errors and system messages. Outputs to File, Database, or Console.
// Log Level is determined by the configuration file
// Controls exit of server if an error occurs after the errors have been written
package logger

import (
  "fmt"
  //"config"
  "os"
)

const (
  thisModule = "Logger"
  typeRequest = "LogType"
  logLevelRequest = "LogLevel"
  sysLevelRequest = "SysLevel"
)

var (
  logType = 3
  logLevel = 3
  sysLevel = 3
)

func SystemMessage(message string, module string) {
  fmt.Println("Saviour::" + module + "::" + message)
}

// Error outputs a assembled error message to the console
func Error(message string, module string, level int) {
  fmt.Println("Error::" + module + "::" + message)
  if (level == 1) {
    os.Exit(1)
  }
}