package core

import "testing"

var ThisUser *User

func TestUser_CheckUserLogin(t *testing.T) {
	var exists bool
	exists, ThisUser = InitUser(InitDatabase(), "Admin", "Password")
	if exists {
		// Pass
	} else {
		t.Error("LoadingAdminFailed")
	}
}
