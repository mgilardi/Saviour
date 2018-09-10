package core

import (
	"testing"
)

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

func TestAccess_InitAccess(t *testing.T) {
	InitOptions()
	InitDebug(true)
	InitCron()
	InitDatabase()
	InitLogger()
	InitCache()
	InitAccess()
	// Commented because testing requires server emulation
	//InitSystem()
}

func TestAccess_createRole(t *testing.T) {
	Logger("TestingRoleCreation", ACCESS, MSG)
	AccessHandler.CreateRole("test")
	if !AccessHandler.CheckRole("test") {
		Logger("TestRoleNotCreated", ACCESS, ERROR)
		t.Error("RoleNotCreated")
	}
}

func TestAccess_cloneRole(t *testing.T) {
	var clonedPerm, testPerm []*Perm
	match := true
	Logger("TestingRoleCloneing", ACCESS, MSG)
	AccessHandler.CloneRole("testclone", "authorized.user")
	perms := AccessHandler.GetPerms()
	for _, perm := range perms {
		for _, role := range perm.Roles {
			if role == "authorized.user" {
				clonedPerm = append(clonedPerm, perm)
			}
			if role == "test" {
				testPerm = append(clonedPerm, perm)
			}
		}
	}
	switch {
	case testPerm.Name == "":
		Logger("FailedToGetTestPerm", ACCESS, ERROR)
		t.Error("FailedToGetTestPerm")
		match = false
	case clonedPerm.Name == "":
		Logger("FailedToGetClonedPerm", ACCESS, ERROR)
		t.Error("FailedToGetClonedPerm")
		match = false
	case len(testPerm.Roles) != len(clonedPerm.Roles):
		Logger("RoleLengthDoesNotEqual", ACCESS, ERROR)
		t.Error("RoleLengthDoesNotEqual")
		match = false
	default:
		for i, role := range testPerm.Roles {
			if clonedPerm.Roles[i] != role {
				Logger("FailedMisMatchedRoleArray", ACCESS, ERROR)
				t.Error("FailedMisMatchedRoleArray")
				match = false
			}
		}
	}
	if !match {
		t.Fatal("FailedRolesDoNotMatch")
	} else {
		t.Log("RolesMatched")
	}
}

func TestAccess_CheckUserAccess(t *testing.T) {
	exists, userMap, user := InitUser("Admin", "Password")
	switch {
	case !exists:
		t.Error("UserDoesNotExist")
	case !AccessHandler.CheckUserAccess(user, "removeuser"):
		t.Error("UserNotAuthorized")
	default:
		t.Log("UserAuthorized")
	}
}

func TestAccess_RemoveRole(t *testing.T) {
	AccessHandler.RemoveRole("test")
	testRoleExists := AccessHandler.CheckRole("test")
	if testRoleExists {
		t.Error("TestRoleNotRemoved")
	}
	AccessHandler.RemoveRole("testclone")
	testCloneRole := AccessHandler.CheckRole("testclone")
	if testCloneRole {
		t.Error("TestCloneRoleNotRemoved")
	}
}
