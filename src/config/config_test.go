package config_test

import (
  "testing"
  "config"
)

func TestGetSettings(t *testing.T) {
  options := config.GetSettings()
  for _, opt := range *options {
    name := opt.FindValue("Name")
    if name == "" {
      t.Errorf("GetSettings Failed")
    }
  }
}

func TestFindValuePass(t *testing.T) {
  value := config.FindValue("Access", "Name")
  if value == nil {
    t.Errorf("Value is NULL")
  }
  if value.(string) == "" {
    t.Errorf("Could Not Get Value")
  }
}

func TestFindValueFail(t *testing.T) {
  value := config.FindValue("Test", "Test")
  if value == nil {
    // good
  } else {
    t.Errorf("FindValue::KeyFound::" + value.(string))
  }
}