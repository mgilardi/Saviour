package core

import (
	"fmt"
	"testing"
)

func TestGetOptions(t *testing.T) {
	InitOptions()
	options := OptionsHandler.GetOptions("Core")
	_, exists := options["Name"]
	if exists && options["Name"].(string) == "Core" {
		fmt.Println("Testing::" + options["Name"].(string))
	} else {
		t.Errorf("CouldNotLoadModule")
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
