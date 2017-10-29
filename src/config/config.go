// Package config loads configuration json files design goals are as follows.
// Must be able to load new values inputted into the configuration without
// modification. Can return a array including all the separate modules with
// their corresponding configuration keys and values.
package config

import (
  "fmt"
  "os"
  "encoding/json"
  "io/ioutil"
  "errors"
)

//Setting stores the id of the module and a map of the imported keys and values
type Setting struct {
  id int
  setMap map[string]string
}

// Constructor For Setting Strut
func initSet(path string, id int) Setting {
  var options Setting
  options.id = id
  options.setMap = make(map[string]string)
  options.setMap["path"] = path
  options.loadFile()
  return options
}

// Load JSON file.
func (options Setting) loadFile() {
  raw, err := ioutil.ReadFile(options.setMap["path"])
  if (err != nil) {
    fmt.Println("Saviour::Config::Error::CannotReadFile")
    // Handle Error
  }
  err = json.Unmarshal(raw, &options.setMap)
  if (err != nil) {
    // Handler Error
  }
}

// Find Configuration Value Using Key
func (options Setting) findValue(key string) string {
  value := options.setMap[key]
  return value
}

// GetSettings loads settings from each modules settings.json into a struct
// array of maps for each module. It returns a pointer to the assemble struct
// array containing the settings map for each module.
func GetSettings() *[]Setting {
  settings := make([]Setting, 0)
  files, err := ioutil.ReadDir("modules")
  if (err != nil) {
    os.Exit(1)
  }
  for i, file := range files {
    settings = append(settings, initSet("modules/" + file.Name() +
      "/settings.json", i))
  }
  return &settings
}

// FileValue returns a setting key from the setting array, if no value is
// found it returns an error.
func FindValue(module string, key string) (error, string) {
  var output string
  var err error
  options := GetSettings()
  for _, opt := range *options {
    if opt.setMap["Name"] == module {
      output = opt.findValue(key)
      break
    }
  }
  if (output == "") {
    err = errors.New("KeyNotFound")
  }
  return err, output
}
