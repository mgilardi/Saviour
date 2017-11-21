package main

import (
	"testing"
)

var db *Database
var uid int

func TestDatabase_InitDatabase(t *testing.T) {
	InitDebug(true)
	InitCron()
	db = InitDatabase()
}

func TestDatabase_CheckUserLogin(t *testing.T) {
	var exists bool
	db.CreateUser("Test", "Test", "Test@test.com")
	exists, uid = db.CheckUserLogin("Test", "Test")
	if exists != true {
		t.Error("CheckUserLogin::FailedUserVerify")
	}
	exists, uid = db.CheckUserLogin("Test", "NotPassword")
	if exists != false {
		t.Error("CheckUserLogin::FailedUserVerify")
	}
	exists, uid = db.CheckUserLogin("NotUser", "NotPassword")
	if exists != false {
		t.Error("CheckUserLogin::FailedUserVerify")
	}
	db.RemoveUser("Test")
}

func TestDatabase_CheckUserExists(t *testing.T) {
	if !db.CheckUserExists("Admin") {
		t.Error("CheckUserExists::Failed::Exists")
	}
	if db.CheckUserExists("NotUser") {
		t.Error("CheckUserExists::Failed::NotExist")
	}
}

func TestDatabase_GetUserMap(t *testing.T) {
	var err error
	uid, err = db.GetUserID("Admin")
	if err != nil {
		t.Error("ErrorGettingIDAdmin::" + err.Error())
	}
	userMap, err := db.GetUserMap(uid)
	if err != nil {
		t.Error("ErrorGettingUserMapAdmin::" + err.Error())
	}
	if !(userMap["Name"].(string) == "Admin") {
		t.Error("")
	}
	_, err = db.GetUserMap(999)
	if err != nil {
		// Correct Outcome
	} else {
		t.Error("NonExistentUserFound")
	}
}

func TestDatabase_Token(t *testing.T) {
	if !db.CheckToken(uid) {
		t.Error("CheckTokenFailed::NotFound")
	}
	if db.CheckToken(999) {
		t.Error("CheckTokenFailed::Found")
	}
	db.CreateUser("Test", "Test", "test@test.com")
	_, uid = db.CheckUserLogin("Test", "Test")
	db.StoreToken(uid, GenToken(32))
	if !db.CheckToken(uid) {
		t.Error("StoreTokenFailed::NotFound")
	}
	db.RemoveUser("Test")
}
