// Package core config loads configuration json files design goals are as follows.
// Must be able to load new values inputted into the configuration without
// modification. Can return a array including all the separate modules with
// their corresponding configuration keys and values.
package core

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	// MODULEOPT is the module name constant for options
	MODULEOPT = "Options"
)

var OptionsHandler Options

type Options struct {
	configPath    string
	loadedOptions map[string]map[string]interface{}
}

// Load Options Struct
func InitOptions() {
	var opt Options
	var err error
	// Find the Config Path
	opt.configPath, err = FindPath("config")
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEOPT+".InitOptions", FATAL)
	}
	// Initialize Config Cache Map
	opt.loadedOptions = make(map[string]map[string]interface{})
	OptionsHandler = opt
}

// GetOption returns either loads the config file or returns the already loaded config
func (opt Options) GetOption(module string) map[string]interface{} {
	optionMap, exists := opt.loadedOptions[module]
	if !exists {
		optionMap = opt.loadOption(module)
	}
	return optionMap
}

// FindValue returns a value of a module
func (opt Options) FindValue(module string, key string) interface{} {
	var output interface{}
	optionsMap := opt.GetOption(module)
	output = optionsMap[key]
	return output
}

// loadOption will read the json config file and load
// it into the config cache and return it to the GetOption function
func (opt Options) loadOption(module string) map[string]interface{} {
	var optionMap map[string]interface{}
	optionMap = make(map[string]interface{})
	optionMap["Path"] = opt.configPath + strings.ToLower(module) + ".json"
	raw, err := ioutil.ReadFile(optionMap["Path"].(string))
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEOPT+".loadOption", ERROR)
	}
	err = json.Unmarshal(raw, &optionMap)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULEOPT+".loadOption", ERROR)
	}
	opt.loadedOptions[module] = optionMap
	return optionMap
}

// FindPath will return the path to any subdirectory of the Saviour folder thats
// within the path of the working directory
func FindPath(searchPath string) (string, error) {
	var assembledPath []string
	Logger("SearchForConfig", PACKAGE+"."+MODULEOPT+".FindConfig", MSG)
	wd, err := os.Getwd()
	// Split the working directory path
	paths := strings.Split(wd, string(os.PathSeparator))
	// Iterate over the array and add the elements until you find the Saviour directory
	for _, folder := range paths {
		assembledPath = append(assembledPath, folder)
		// If Saviour subdirectory is found iterate that directory for the specified search path
		if folder == "Saviour" {
			Logger("ConfigFoundSaviourFolder", PACKAGE+"."+MODULEOPT+".FindConfig", MSG)
			Logger(path.Join(assembledPath...), PACKAGE+"."+MODULEOPT+".FindConfig", MSG)
			dirs, _ := ioutil.ReadDir("/" + path.Join(assembledPath...))
			if err != nil {
				Logger("CouldNotOpenAssembledPath::"+path.Join(assembledPath...), PACKAGE+"."+MODULEOPT+".FindConfig", MSG)
			}
			// If searchPath directory is found within the Saviour directory generate a path
			for _, dir := range dirs {
				if dir.Name() == searchPath {
					Logger("ConfigFolderFound", PACKAGE+"."+MODULEOPT+".FindConfig", MSG)
					assembledPath = append(assembledPath, dir.Name())
					return "/" + path.Join(assembledPath...) + "/", err
				}
			}
		}
	}
	// If the directory is not found return an error
	return "Error", errors.New("PathFailedToLocateDirectory")
}
