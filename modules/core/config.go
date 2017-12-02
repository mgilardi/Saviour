// Package core config loads configuration json files design goals are as follows.
// Must be able to load new values inputted into the configuration without
// modification. Can return a array including all the separate modules with
// their corresponding configuration keys and values.
package core

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

// OptionsHandler global variable for options struct
var OptionsHandler *Options

// Option Struct
type Option struct {
	name string
}

// Create a new option object
func initOption(name string) *Option {
	var option Option
	option.name = name
	return &option
}

// GetValues for the option object from cache
func (opt *Option) GetValues() map[string]interface{} {
	exists, cacheMap := CacheHandler.Cache(opt)
	if !exists {
		Logger("CouldNotLoadCacheOptions::"+opt.name, "Config", ERROR)
	}
	return cacheMap
}

// UpdateCache updates the cache entry for this option object
func (opt *Option) UpdateCache() {
	CacheHandler.Update(opt)
}

// Cache This function is triggered by the cache module
func (opt *Option) Cache() (string, map[string]interface{}) {
	cacheID := opt.name + ":option"
	values := opt.loadValues()
	return cacheID, values
}

// This loads the values for the options object from disk
func (opt *Option) loadValues() map[string]interface{} {
	options := getOptions(opt.name)
	return options
}

// CacheID This is a helper function for the cache module
func (opt *Option) CacheID() string {
	cacheID := opt.name + ":option"
	return cacheID
}

// Options struct
type Options struct {
	options     map[string]*Option
	cacheLoaded bool
}

// InitOptions initializes the options struct
// @TODO Get path dynamically instead of hard coding: /src/Saviour/modules
func InitOptions() {
	var opt Options
	opt.cacheLoaded = false
	opt.options = make(map[string]*Option)
	dir, err := ioutil.ReadDir(os.Getenv("GOPATH") + "/src/Saviour/modules")
	if err != nil {
		Logger(err.Error(), "Config", ERROR)
	}
	for _, file := range dir {
		Logger("LoadingOptionsFile::"+file.Name(), "Config", MSG)
		newOption := initOption(file.Name())
		opt.options[file.Name()] = newOption
	}
	OptionsHandler = &opt
}

// GetOptions will return the options map
func (opts *Options) GetOptions(module string) map[string]interface{} {
	if opts.cacheLoaded == false {
		return getOptions(module)
	}
	return opts.options[module].GetValues()
}

// CacheLoaded is the switch from loading from disk to loading from cache
func (opts *Options) CacheLoaded() {
	opts.cacheLoaded = true
}

// GetOptions returns an map with the loaded options from the json settings file
// @TODO Find the Go equivellent of: PHP's __FUNCTION__ . __FILE__ . __LINE__
func getOptions(module string) map[string]interface{} {
	var optionsMap map[string]interface{}
	optionsMap = make(map[string]interface{})
	optionsMap["Path"] = os.Getenv("GOPATH") + "/src/Saviour/modules/" +
		strings.ToLower(module) + "/settings.json"
	raw, err := ioutil.ReadFile(optionsMap["Path"].(string))
	if err != nil {
		Logger(err.Error(), "Config", ERROR)
	}
	err = json.Unmarshal(raw, &optionsMap)
	if err != nil {
		Logger(err.Error(), "Config", ERROR)
	}
	return optionsMap
}

// FindValue returns a value of a module
func FindValue(module string, key string) interface{} {
	var output interface{}
	optionsMap := getOptions(module)
	output = optionsMap[key]
	return output
}
