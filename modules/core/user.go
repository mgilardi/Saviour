package core

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"
)

// User handles users
type User struct {
	uid                int
	name, token, email string
	db                 *Database
	online             bool
}

// InitUser constructs the user on initial login
// CheckUserLogin inputs username and password and checks if it matches
// the database returns a boolean if the user if found along with the
// users database uid
func InitUser(name string, pass string) (bool, *User) {
	var dbPass string
	var uid int
	var verified = false
	var user User
	user.online = false
	user.db = DBHandler
	Sys("CheckUserLogin::"+name, "User")
	err := user.db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		Sys("UserNotFound", "User")
	case err != nil:
		Error(err, "User")
	case dbPass == pass:
		Sys("LoginVerified::"+name, "User")
		user.name = name
		user.uid = uid
		user.CheckTokenExists()
		verified = true
	default:
		Warn(errors.New("InvalidPassword::"+name), "User")
	}
	return verified, &user
}

// CheckTokenExists checks to see if a token exists in the database if not
// it generates one.
func (user *User) CheckTokenExists() {
	exists, _ := user.CheckToken(user.uid)
	if !exists {
		newToken := GenToken(32)
		user.StoreToken(user.uid, newToken)
	}
}

// VerifyToken verify the user token
func (user *User) VerifyToken(token string) bool {
	if token == user.GetToken() {
		return true
	}
	return false
}

// GetName returns username
func (user *User) GetName() string {
	return user.name
}

// GetEmail returns email from cache
func (user *User) GetEmail() string {
	cacheMap := user.GetCache()
	return cacheMap["email"].(string)
}

// GetToken returns user token
func (user *User) GetToken() string {
	cacheMap := user.GetCache()
	return cacheMap["token"].(string)
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
	err := user.db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		Sys("TokenNotFound", "User")
	case err != nil:
		Error(err, "User")
	default:
		Sys("TokenFound", "User")
		exists = true
	}
	return exists, token
}

// StoreToken writes user token to the database
func (user *User) StoreToken(uid int, token string) {
	Sys("StoreNewToken", "User")
	_, err := user.db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE token = ?`, uid, token, token)
	if err != nil {
		Error(err, "User")
	}
}

// GetCache requests user info from cache
func (user *User) GetCache() map[string]interface{} {
	exists, userCache := CacheHandler.Cache(user)
	if !exists {
		Error(errors.New("FailedToLoadCache"), "User")
	}
	return userCache
}

// UpdateCache updates user cache entrys
func (user *User) UpdateCache() {
	CacheHandler.Update(user)
}

// Cache is called to load user information into the cache
func (user *User) Cache() (string, map[string]interface{}) {
	cid := user.GetName() + ":" + strconv.Itoa(user.uid) + ":UserData"
	cacheMap := user.loadCacheValues()
	return cid, cacheMap
}

// CacheID returns the cache id for this user
func (user *User) CacheID() string {
	return user.GetName() + ":" + strconv.Itoa(user.uid) + ":UserData"
}

// loadCacheValues loads the user information from the database and places
// them in a map for loading into the cache
func (user *User) loadCacheValues() map[string]interface{} {
	var uid int
	var name, email, token string
	cacheMap := make(map[string]interface{})
	user.db.sql.QueryRow(`SELECT users.uid, name, mail, token FROM users JOIN login_token WHERE users.uid = login_token.uid & users.uid = ?`, user.uid).Scan(&uid, &name, &email, &token)
	cacheMap["uid"] = uid
	cacheMap["name"] = name
	cacheMap["email"] = email
	cacheMap["token"] = token
	return cacheMap
}

// GenToken generates user token of specified length
func GenToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		Error(err, "User")
	}
	return base64.URLEncoding.EncodeToString(byte)
}

// GetUserID will return the database uid for a username
func GetUserID(db *Database, name string) (int, error) {
	var uid int
	Sys("GetUserID::"+name, "User")
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		Sys("UserNotFound", "User")
	case err != nil:
		Error(err, "User")
	}
	return uid, err
}

// CheckUserExists checks the database for a username and returns true or false
// if it exists
func CheckUserExists(db *Database, name string) bool {
	Sys("CheckUserExists::"+name, "User")
	exists := false
	_, err := GetUserID(db, name)
	switch {
	case err == sql.ErrNoRows:
		Sys("UserNotFound", "User")
	case err != nil:
		Error(err, "User")
	default:
		exists = true
	}
	return exists
}
