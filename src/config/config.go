/*

Saviour Configuration Handler

Provide Unified Way To Provide Pre-Defined Settings To Modules

Load JSON file settings.json from each module

Provide detection of changes in settings.json and load changes to cache 

*/

package config

import ("fmt"
		"os"
		"encoding/json"
		"io/ioutil")

// Setting Stuct
type Setting struct {
	id int
	setMap map[string]string
} 

// Constructor For Setting Strut
func initSet(path string, i int) Setting {
	s:= Setting{id: i}
	s.setMap = make(map[string]string)
	s.setMap["path"] = path
	s.loadFile()
	return s
}


// Load JSON file.
func (s Setting) loadFile() {
	raw, err := ioutil.ReadFile(s.setMap["path"])
	if err != nil {
		fmt.Println("Saviour::Config::Error::CannotReadFile")
		// Handle Error
	}
	err = json.Unmarshal(raw, &s.setMap)
	if err != nil {
		// Handler Error
	}
}

// Find Configuration Value Using Key
func (s Setting) findValue(key string) string {
	value := s.setMap[key]
	return value
}

/*

GetSettings() Function as described in the design document:

Load Settings From JSON Files into a struct array

Settings Handler must be indepented for each module and not require 
modification when new modules are added.

This is accomplished by using a map to store settings. To 
set module settings each key and value will be taken from the json 
file.

When an outside source needs the configuration value it will search
the Settings Array for the module and then request a specific key
which will return the value.

*/	
func GetSettings() []Setting {
	// Move to a status/log messaging service
	fmt.Println("Saviour::Config::Start")
	fmt.Println("Saviour::Config::LoadModules")
	settings := make([]Setting, 0)
	files, err := ioutil.ReadDir("modules")
	if err != nil {
		os.Exit(1)
	}
	i := 0
	for {
		if (i == len(files)) {
			break
		} else {
			file := files[i] 
			settings = append(settings, initSet("modules/" + file.Name() + "/settings.json", i))
			fmt.Println("Saviour::Config::Module::" + settings[i].setMap["Name"])
			i++
		}
	}
	return settings
}

func FindValue(module string, key string, s []Setting) string {
	i := 0
	var output string
	for {
		if (i == len(s)) {
			break
		} 
		if (s[i].setMap["Name"] == module) {
			output = s[i].setMap[key]
			break
		}
		i++
	}
	return output
}