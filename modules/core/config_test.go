package core

import (
	"fmt"
	"testing"
)

func TestGetOptions(t *testing.T) {
	options := GetOptions("Core")
	if options["Name"].(string) == "Core" {
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
	value := FindValue("Core", "Name")
	if value == nil {
		t.Errorf("ValueIsNull")
	}
	if value.(string) == "" {
		t.Errorf("Could Not Get Value")
	}
}
