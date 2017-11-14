PACKAGE DOCUMENTATION

package config
    import "config"

    Package config loads configuration json files design goals are as
    follows. Must be able to load new values inputted into the configuration
    without modification. Can return a array including all the separate
    modules with their corresponding configuration keys and values.

FUNCTIONS

func FindValue(module string, key string) interface{}
    FindValue returns a value of a module

func GetAllOptions() []map[string]interface{}
    GetAllOptions returns a array of maps with the loaded settings.json
    files for each module

func GetOptions(module string) map[string]interface{}
    GetOptions returns an map with the loaded options from the json settings
    file
