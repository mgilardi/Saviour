package core

import (
	"testing"
)

func TestCommand_CreateUser(t *testing.T) {
	err := CommandHandler.CreateUser("TestCommand", "Password", "test@command.com")
	if err != nil {
		t.Error("CreateUserFailed")
	}
	exists, _ := GetUserID("TestCommand")
	if !exists {
		t.Error("TestCommandUserNotCreated")
	}
	err = CommandHandler.CreateUser("", "Password", "Email")
	if err == nil {
		t.Error("CreateConditionFailed::" + err.Error())
	}
	err = CommandHandler.CreateUser("TestCommand", "", "Email")
	if err == nil {
		t.Error("CreateConditionFailed::" + err.Error())
	}
	err = CommandHandler.CreateUser("TestCommand", "Password", "")
	if err == nil {
		t.Error("CreateConditionFailed::" + err.Error())
	}
}

func TestCommand_RemoveUser(t *testing.T) {
	_, _, currentUser := InitUser("Admin", "Password")
	output := CommandHandler.RemoveUser("TestCommand", currentUser)
	t.Log("Output::" + output)
}
