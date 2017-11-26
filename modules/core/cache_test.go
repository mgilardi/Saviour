package core

import (
	"fmt"
	"testing"
)

var testHandler *test

type test struct {
	//
}

func initTest() {
	var t test
	testHandler = &t
}

func (t test) Cache() (string, map[string]interface{}) {
	name := "test"
	cacheMap := make(map[string]interface{})
	cacheMap["test"] = "ThisIsATest"
	return name, cacheMap
}

func (t test) CacheID() string {
	return "test"
}

func TestCache_InitCache(t *testing.T) {
	InitDebug(true)
	InitOptions()
	InitCron()
	InitDatabase()
	InitLogger()
	InitCache()
	InitSystem()
	initTest()
	exists, testCache := CacheHandler.GetCache(testHandler)
	if !exists {
		t.Error("GetCacheFailed")
	}
	fmt.Println("Output:" + testCache["test"].(string))
	CacheHandler.CheckCache()
	CacheHandler.ClearCache()
}
