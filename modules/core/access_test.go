package core

import "testing"

var testHandler Test

type Test struct {
	//
}

func initTest() Test {
	var test Test
	return test
}

func (t Test) Cache() (string, map[string]interface{}) {
	name := "test"
	cacheMap := make(map[string]interface{})
	cacheMap["test"] = "ThisIsATest"
	return name, cacheMap
}

func (t Test) CacheID() string {
	return "test"
}

func TestCache_InitAccess(t *testing.T) {
	InitOptions()
	InitDebug(true)
	InitCron()
	InitDatabase()
	InitLogger()
	InitCache()
	InitAccess()
	//InitSystem()
}
