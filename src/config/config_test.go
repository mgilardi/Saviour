package config

import (
  "testing"
  "io/ioutil"
  "reflect"
)

func TestGetSettings(t *testing.T) {
  var options *[]Setting
  var optionsTest []Setting
  var optionsTestPtr *[]Setting
  options = GetSettings()
  files, _ := ioutil.ReadDir("modules")
  for i, file := range files {
    optionsTest = append(optionsTest, initSet("modules/" + file.Name() +
      "/settings.json", i))
  }
  optionsTestPtr = &optionsTest
  if (!reflect.DeepEqual(options, optionsTestPtr)) {
    t.Errorf("TestGetSettings:OutputsDontMatch")
  }
}

func TestFindValuePass(t *testing.T) {
  err, _ := FindValue("Access", "Name")
  if err != nil {
    t.Errorf("FindValue::KeyNotFound")
  }
}

func TestFindValueFail(t *testing.T) {
  err , value := FindValue("Test", "Test")
  if (err != nil) {
    // No Fail
  }
  if (value != "") {
    t.Error("FindValue::KeyFound")
  }
}