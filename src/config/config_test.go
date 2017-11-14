package config_test

import (
	"config"
	"fmt"
	"testing"
)

func TestGetOptions(t *testing.T) {
	options := config.GetOptions("Database")
	if options["Name"].(string) == "Database" {
		fmt.Println("Testing::" + options["Name"].(string))
	} else {
		t.Errorf("CouldNotLoadModule")
	}
}

func TestGetAllOptions(t *testing.T) {
	optionsArray := config.GetAllOptions()
	for _, opt := range optionsArray {
		if opt["Name"].(string) == "" {
			t.Errorf("CouldNotLoadModule")
		}
	}
}

func TestFindValuePass(t *testing.T) {
	value := config.FindValue("Access", "Name")
	if value == nil {
		t.Errorf("ValueIsNull")
	}
	if value.(string) == "" {
		t.Errorf("Could Not Get Value")
	}
}
