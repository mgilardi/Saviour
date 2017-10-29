#PACKAGE DOCUMENTATION

package config
    import "config"

    Package config loads configuration json files design goals are as
    follows. Must be able to load new values inputted into the configuration
    without modification. Can return a array including all the separate
    modules with their corresponding configuration keys and values.

##FUNCTIONS

func FindValue(module string, key string) (error, string)
    FileValue returns a setting key from the setting array, if no value is
    found it returns an error.

func GetSettings() *[]Setting
    GetSettings loads settings from each modules settings.json into a struct
    array of maps for each module. It returns a pointer to the assemble
    struct array containing the settings map for each module.

##TYPES

type Setting struct {
    // contains filtered or unexported fields
}
    Setting stores the id of the module and a map of the imported keys and
    values