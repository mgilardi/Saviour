package core

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
)

// User handles users
type User struct {
	uid         int
	name, token string
	db          *Database
	cache       *Cache
	online      bool
}

// InitUser constructs the user on initial login
// CheckUserLogin inputs username and password and checks if it matches
// the database returns a boolean if the user if found along with the
// users database uid
func InitUser(db *Database, name string, pass string) (bool, *User) {
	var dbPass string
	var uid int
	var verified = false
	var user User
	user.online = false
	user.db = db
	DebugHandler.Sys("CheckUserLogin::"+name, "User")
	err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", "User")
	case err != nil:
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 3)
	case dbPass == pass:
		DebugHandler.Sys("LoginVerified::"+name, "User")
		verified = true
		user.uid = uid
		user.name = name
		user.CheckTokenExists()
	default:
		DebugHandler.Sys("InvalidPassword::"+name, "User")
		LogHandler.Warn(errors.New("InvalidPassword::"+name), "User")
	}
	user.cache = CacheHandler

	return verified, &user
}

// CheckTokenExists checks to see if a token exists in the database if not
// it generates one.
func (user *User) CheckTokenExists() {
	exists, token := user.CheckToken(user.uid)
	if !exists {
		user.token = GenToken(32)
		user.StoreToken(user.uid, user.token)
	} else {
		user.token = token
	}
}

// GetName returns username
func (user *User) GetName() string {
	return user.name
}

// GetToken returns user token
func (user *User) GetToken() string {
	return user.token
}

// SetToken generates a new token and writes it to DB
func (user *User) SetToken() {
	user.token = GenToken(32)
	user.StoreToken(user.uid, user.token)
}

// IsOnline returns the online flag for the user
func (user *User) IsOnline() bool {
	return user.online
}

// SetOnline will set the flag to the input
func (user *User) SetOnline(isOnline bool) {
	user.online = isOnline
}

// CheckToken checks if the user has a token in the database login_token table
func (user *User) CheckToken(uid int) (bool, string) {
	var token string
	exists := false
	DebugHandler.Sys("CheckingToken", "User")
	err := user.db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("TokenNotFound", "User")
	case err != nil:
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 3)
	default:
		DebugHandler.Sys("TokenFound", "User")
		exists = true
	}
	return exists, token
}

// StoreToken writes user token to the database
func (user *User) StoreToken(uid int, token string) {
	DebugHandler.Sys("StoreToken", "User")
	_, err := user.db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE token = ?`, uid, token, token)
	if err != nil {
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 3)
	}
}

// GenToken generates user token of specified length
func GenToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 1)
	}
	return base64.URLEncoding.EncodeToString(byte)
}

// GetUserID will return the database uid for a username
func GetUserID(db *Database, name string) (int, error) {
	var uid int
	DebugHandler.Sys("GetUserID::"+name, "User")
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", "User")
	case err != nil:
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 3)
	}
	return uid, err
}

// CheckUserExists checks the database for a username and returns true or false
// if it exists
func CheckUserExists(db *Database, name string) bool {
	DebugHandler.Sys("CheckUserExists::"+name, "User")
	exists := false
	_, err := GetUserID(db, name)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", "User")
	case err != nil:
		LogHandler.Err(err, "User")
		DebugHandler.Err(err, "User", 3)
	default:
		exists = true
	}
	return exists
}
