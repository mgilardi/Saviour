package main

import (
	"fmt"
	"testing"
)

func TestGetOptions(t *testing.T) {
	options := GetOptions("Database")
	if options["Name"].(string) == "Database" {
		fmt.Println("Testing::" + options["Name"].(string))
	} else {
		t.Errorf("CouldNotLoadModule")
	}
}

func TestGetAllOptions(t *testing.T) {
	optionsArray := GetAllOptions()
	for _, opt := range optionsArray {
		if opt["Name"].(string) == "" {
			t.Errorf("CouldNotLoadModule")
		}
	}
}

func TestFindValuePass(t *testing.T) {
	value := FindValue("Access", "Name")
	if value == nil {
		t.Errorf("ValueIsNull")
	}
	if value.(string) == "" {
		t.Errorf("Could Not Get Value")
	}
}
