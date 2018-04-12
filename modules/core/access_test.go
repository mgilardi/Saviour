package core

import "testing"

var testHandler Test

type Test struct {
	//
}

func initTest() Test {
	var test Test
	AccessHandler.LoadPerm(test)
	return test
}

func (t Test) PermID() string {
	return "testPerm"
}

func (t Test) DefaultPerm() map[string]bool {
	accessMap := make(map[string]bool)
	accessMap["admin"] = true
	return accessMap
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
	testThis := initTest()
	_, _, currentUser := InitUser("Admin", "Password")
	pass := AccessHandler.CheckPerm(currentUser, testThis)
	if !pass {
		t.Error("AccessFailed")
	}
	AccessHandler.DisableAccess("admin", testThis)
	pass = AccessHandler.CheckPerm(currentUser, testThis)
	if pass {
		t.Error("AccessFailed::AdminAccessDisabled")
	}
	AccessHandler.AllowAccess("admin", testThis)
	AccessHandler.clearDB()
}
