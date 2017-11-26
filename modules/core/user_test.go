package core

import (
	"strconv"
	"testing"
)

var ThisUser *User

func TestUser_CheckUserLogin(t *testing.T) {
	var exists bool
	exists, ThisUser = InitUser("Admin", "Password")
	if exists {
		// Pass
		Sys("Loading::"+ThisUser.GetName(), "Test")
		Sys("Token::"+ThisUser.GetToken(), "Test")
		Sys("IsOnline::"+strconv.FormatBool(ThisUser.IsOnline()), "Test")
	} else {
		t.Error("LoadingAdminFailed")
	}
}
