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

func TestUser_GetUser(t *testing.T) {
	exists, userMap, currentUser := InitUser("Admin", "Password")
	if exists {
		getUserExists, _, _ := GetUser(userMap["name"].(string), currentUser.GetToken())
		if !getUserExists {
			t.Error("GetUserFailed")
		}
	}
}

func TestUser_GetUnAuthorizedUser(t *testing.T) {
	unAuthUser := GetUnAuthorizedUser()
	if unAuthUser.GetName() != "UnAuthorized" {
		t.Error("GetUnAuthorizedUserFailed")
	}
}

func TestUser_VerifyToken(t *testing.T) {
	var verify bool
	exists, _, currentUser := InitUser("Admin", "Password")
	if exists {
		verify = currentUser.VerifyToken(currentUser.GetToken())
	}
	if !verify {
		t.Error("VerifyTokenFailed")
	}
}

func TestUser_GetRoleNames(t *testing.T) {
	var verified bool
	exists, _, currentUser := InitUser("Admin", "Password")
	if !exists {
		t.Error("InitUserFailed")
	} else {
		roleNames := currentUser.GetRoleNames()
		for role, _ := range roleNames {
			if role == "administrator" {
				verified = true
				break
			} else {
				verified = false
			}
		}
		if !verified {
			t.Error("RoleNotAdministrator")
		}
	}
}

func TestUser_GetEmail(t *testing.T) {
	exists, _, currentUser := InitUser("Admin", "Password")
	if !exists {
		t.Error("InitUserFailed")
	} else {
		email := currentUser.GetEmail()
		if email != "ian@diysecurity.com" {
			t.Error("GetEmailMisMatch")
		}
	}
}

func TestUser_SetPassword(t *testing.T) {
	exists, _, currentUser := InitUser("Admin", "Password")
	if !exists {
		t.Error("AdminUserNotFound")
	} else {
		currentUser.SetPassword("Pass")
		exists, _, currentUser = InitUser("Admin", "Pass")
		if !exists {
			t.Error("ChangePasswordFailed")
		} else {
			currentUser.SetPassword("Password")
		}
	}
}

func TestUser_StoreToken(t *testing.T) {
	token := GenToken(32)
	exists, _, currentUser := InitUser("Admin", "Password")
	if !exists {
		t.Error("AdminUserNotFound")
	} else {
		currentUser.StoreToken(currentUser.uid, token)
	}
}
