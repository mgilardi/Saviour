// Loads configuration json files design goals are as follows.
// Must be able to load new values inputted into the configuration without
// modification. Can return a array including all the separate modules with
// their corresponding configuration keys and values.
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const (
	thisModuleConfig = "Config"
)

// GetOptions returns an map with the loaded options from the json settings file
func GetOptions(module string) map[string]interface{} {
	optionsMap := make(map[string]interface{})
	optionsMap["Path"] = os.Getenv("GOPATH") + "/src/Saviour/modules/" +
		strings.ToLower(module) + "/settings.json"
	raw, err := ioutil.ReadFile(optionsMap["Path"].(string))
	if err != nil {
		DebugHandler.Err(err, thisModuleConfig, 1)
	}
	err = json.Unmarshal(raw, &optionsMap)
	if err != nil {
		DebugHandler.Err(err, thisModuleConfig, 1)
	}
	return optionsMap
}

// GetAllOptions returns a array of maps with the loaded settings.json files for each module
func GetAllOptions() []map[string]interface{} {
	var optionArray []map[string]interface{}
	var currentModule map[string]interface{}
	dir, err := ioutil.ReadDir(os.Getenv("GOPATH") + "/src/Saviour/modules/")
	if err != nil {
		DebugHandler.Err(err, thisModuleConfig, 1)
	}
	for _, file := range dir {
		currentModule = GetOptions(file.Name())
		optionArray = append(optionArray, currentModule)
	}
	return optionArray
}

// FindValue returns a value of a module
func FindValue(module string, key string) interface{} {
	var output interface{}
	optionsMap := GetOptions(module)
	output = optionsMap[key]
	return output
}
