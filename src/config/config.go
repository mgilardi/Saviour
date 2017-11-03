// Package config loads configuration json files design goals are as follows.
// Must be able to load new values inputted into the configuration without
// modification. Can return a array including all the separate modules with
// their corresponding configuration keys and values.
package config

import (
        "os"
        "encoding/json"
        "io/ioutil"
        "errors"
        "fmt"
)

const (
  thisModule = "Config"
)

//Setting stores the id of the module and a map of the imported keys and values
type Setting struct {
  id int
  setMap map[string]interface{}
}

// initSetting is a constructor for the config object initSetting(pathtofile, idofmodule)
// returns an array of config objects for each module.
func initSetting(path string, id int) (error, Setting) {
  var options Setting
  options.id = id
  options.setMap = make(map[string]interface{})
  options.setMap["path"] = path
  err := options.LoadFile()
  return err, options
}

// LoadFile loads the JSON settings file for this specific module.
func (options *Setting) LoadFile() error {
  raw, err := ioutil.ReadFile(options.setMap["path"].(string))
  err = json.Unmarshal(raw, &options.setMap)
  return err
}

// FindValue returns a value from the setting map
func (options *Setting) FindValue(key string) interface{} {
  value := options.setMap[key]
  return value
}

// GetSettings loads settings from each module creates an array of config objects
// with a setting map for each module. It returns a pointer to the assembled object
// array.
func GetSettings() *[]Setting {
  settings := make([]Setting, 0)
  files, err := ioutil.ReadDir("modules")
  if (err != nil) {
    fmt.Println("Error::Config::CannotReadDirectory")
    fmt.Println("Error::Config::" + err.Error())
    os.Exit(1)
  }
  for i, file := range files {
    err, setting := initSetting("modules/" + file.Name() + "/settings.json", i)
    if (err != nil) {
     fmt.Println("Error::Config::CannotLoadSettings")
     fmt.Println("Error::Config::" + err.Error())
    } else {
      settings = append(settings, setting)
    }
  }
  return &settings
}

// GetSettingModule takes in the array of Setting and gives back the element for the
// specified module.
func GetSettingModule(module string, options *[]Setting) (error, *Setting) {
  var err error
  var set Setting
  var check string
  for _, opt := range *options {
    if opt.setMap["Name"] == module {
      set = opt
      check = opt.setMap["Name"].(string)
      break
    }
  }
  if (check == "") {
    err = errors.New("ModuleNotFound")
  }
  return err, &set
}

// FileValue returns a setting key from the setting array
func FindValue(module string, key string) interface{} {
  var output interface{}
  options := GetSettings()
  for _, opt := range *options {
    if (opt.setMap["Name"] == module) {
      output = opt.FindValue(key)
      break
    }
  }
  return output
}
