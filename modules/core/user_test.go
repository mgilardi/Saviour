package core

import (
	"testing"
)

func TestUser_CheckUserLogin(t *testing.T) {
	exists, userMap, currentUser := InitUser("Admin", "Password")
	if exists {
		Logger("Loading::"+currentUser.GetName(), "Test", MSG)
		Logger("Email::"+userMap["email"].(string), "Test", MSG)
		Logger("Token::"+currentUser.GetToken(), "Test", MSG)
	} else {
		t.Error("LoadingUserFailed")
	}
}
