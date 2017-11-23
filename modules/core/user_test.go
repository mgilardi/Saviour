package core

import (
	"strconv"
	"testing"
)

var ThisUser *User

func TestUser_CheckUserLogin(t *testing.T) {
	var exists bool
	exists, ThisUser = InitUser(InitDatabase(), "Admin", "Password")
	if exists {
		// Pass
		DebugHandler.Sys("Loading::"+ThisUser.GetName(), "Test")
		DebugHandler.Sys("Token::"+ThisUser.GetToken(), "Test")
		DebugHandler.Sys("IsOnline::"+strconv.FormatBool(ThisUser.IsOnline()), "Test")
	} else {
		t.Error("LoadingAdminFailed")
	}
}
