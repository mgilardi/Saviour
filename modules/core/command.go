package core

import (
	"database/sql"
	"errors"
	"time"
)

// Admin command contains all adminitstrative commands for the server
var CommandHandler Command

type Command struct {
}

func InitCommand() {
	var cmd Command
	CommandHandler = cmd
}

func (cmd Command) CreateUser(name string, pass string, email string) error {
	var userCheck sql.NullString
	var err error
	Logger("CreateUser::"+name, SYSTEM, MSG)
	DBHandler.sql.QueryRow(`SELECT name FROM users WHERE name = ?`, name).Scan(&userCheck)
	switch {
	case userCheck.Valid:
		Logger("DuplicateUser::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("DuplicateUser::UserCreationFailed")
	case name == "":
		Logger("NameEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("NameEntryIsEmpty::UserCreationFailed")
	case pass == "":
		Logger("PasswordEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("PasswordEntryIsEmpty::UserCreationFailed")
	case email == "":
		Logger("EmailEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("EmailEntryIsEmpty::UserCreationFailed")
	default:
		Logger("CreatingUser::Name:"+name+"::"+email, SYSTEM, MSG)
		hashPass := GenHashPassword(pass)
		currentTime := time.Now().Unix()
		insertUser := DBHandler.SetupExec(
			`INSERT INTO users (name, pass, mail, created) `+
				`VALUES (?, ?, ?, ?)`, name, hashPass, email, currentTime)
		DBHandler.Exec(insertUser)
		_, uid := GetUserID(name)
		insertUserRole := DBHandler.SetupExec(
			`INSERT INTO user_roles (uid, rid) `+
				`VALUES (?, ?)`, uid, 3)
		DBHandler.Exec(insertUserRole)
	}
	return err
}

func (cmd Command) ChangeUserRole(roleChangeUser string, role string, requestUser *User) string {
	Logger("ChangeUserRole::"+roleChangeUser, SYSTEM, MSG)
	//exists, roleChangeUserID := GetUserID(roleChangeUser)
	return "OperationCompleted::User::" + roleChangeUser + "::RoleChange::" + role
}

// RemoveUser removes a user entry from the database
func (cmd Command) RemoveUser(removeUser string, requestUser *User) string {
	Logger("RemoveUser::"+removeUser, SYSTEM, MSG)
	exists, uid := GetUserID(removeUser)
	switch {
	case !exists:
		Logger("CouldNotRemoveUser::DoesNotExist", SYSTEM, ERROR)
		return "UserNotFound"
	default:
		deleteUserToken := DBHandler.SetupExec(`DELETE FROM login_token WHERE uid = ?`, uid)
		deleteUserRoles := DBHandler.SetupExec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		deleteUser := DBHandler.SetupExec(`DELETE FROM users WHERE uid = ?`, uid)
		DBHandler.Exec(deleteUserToken, deleteUserRoles, deleteUser)
		return removeUser + "::Removed"
	}
}
