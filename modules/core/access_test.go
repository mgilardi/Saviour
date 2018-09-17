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
	InitCommand()
	// Commented because testing requires server emulation
	//InitSystem()
}

func TestAccess_CreateRole(t *testing.T) {
	Logger("TestingRoleCreation", ACCESS, MSG)
	AccessHandler.CreateRole("test")
	if !AccessHandler.CheckRole("test") {
		Logger("TestRoleNotCreated", ACCESS, ERROR)
		t.Error("RoleNotCreated")
	}
}

func TestAccess_CloneRole(t *testing.T) {
	var clonedPerm, testPerm []*Perm
	var match bool
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
	for _, perm := range testPerm {
		for _, cloned := range clonedPerm {
			if perm.Name == cloned.Name {
				match = true
			} else {
				match = false
				break
			}
		}
		if !match {
			t.Error("PermissionMisMatch")
			break
		}
	}
}

func TestAccess_CheckUserAccess(t *testing.T) {
	exists, _, user := InitUser("Admin", "Password")
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
