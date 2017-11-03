package config

import (
  "testing"
  "io/ioutil"
  "reflect"
)

func TestGetSettings(t *testing.T) {
  var options []Setting
  var optionsTest Setting
  var optionsTestPtr *[]Setting
  var err error
  options = GetSettings()
  files, _ := ioutil.ReadDir("modules")
  for i, file := range files {
    err, optionsTest = initSetting("modules/" + file.Name() + "/settings.json", i)
    options = append(options, optionsTest)
  }
  if err != nil {

  }
  optionsTestPtr = &optionsTest
  if (!reflect.DeepEqual(options, optionsTestPtr)) {
    t.Errorf("TestGetSettings:OutputsDontMatch")
  }
}

func TestFindValuePass(t *testing.T) {
  value := FindValue("Access", "Name").(string)

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