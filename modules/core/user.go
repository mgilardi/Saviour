package core

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"
)

const (
	thisModuleUser = "User"
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
	user.db = db
	DebugHandler.Sys("CheckUserLogin::"+name, thisModuleDB)
	err := db.sql.QueryRow(`SELECT pass, uid FROM users WHERE name = ?`, name).Scan(&dbPass, &uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	case dbPass == pass:
		verified = true
	default:
		DebugHandler.Sys("InvalidPassword::"+name, thisModuleDB)
		LogHandler.Warn(errors.New("InvalidPassword::"+name), thisModuleDB)
	}
	user.cache = CacheHandler
	user.CheckTokenExists()
	user.InfoUpdate()
	return verified, &user
}

// CheckToken checks to see if a token exists in the database if not
// it generates one.
func (user *User) CheckTokenExists() {
	if !user.CheckToken(user.uid) {
		user.StoreToken(user.uid, GenToken(32))
	}
}

// InfoUpdate calls for the username and token from the database to be
// held in memory
func (user *User) InfoUpdate() {
	userInfo := user.GetInfoMap()
	user.name = userInfo["Name"].(string)
	user.token = userInfo["Token"].(string)
}

// GetInfoMap retrieves the loaded infomap from the database
func (user *User) GetInfoMap() map[string]interface{} {
	exists, userInfo := user.cache.GetCacheMap("user:" + strconv.Itoa(user.uid) + ":" + "info")
	if !exists {
		var err error
		userInfo, err = user.GetUserMap(user.uid)
		if err != nil {
			LogHandler.Err(err, thisModuleUser)
			DebugHandler.Err(err, thisModuleUser, 1)
		}
		user.cache.SetCacheMap("user:"+strconv.Itoa(user.uid)+":"+"info", userInfo, true)
	}
	return userInfo
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
	token := GenToken(32)
	user.StoreToken(user.uid, token)
	user.InfoUpdate()
}

// IsOnline returns the online flag for the user
func (user *User) IsOnline() bool {
	return user.online
}

// SetOnline will set the flag to the input
func (user *User) SetOnline(isOnline bool) {
	user.online = isOnline
}

// GenToken generates user token of specified length
func GenToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		LogHandler.Err(err, thisModuleUser)
		DebugHandler.Err(err, thisModuleUser, 1)
	}
	return base64.URLEncoding.EncodeToString(byte)
}

// CheckUserExists checks the database for a username and returns true or false
// if it exists
func (usr *User) CheckUserExists(name string) bool {
	DebugHandler.Sys("CheckUserExists::"+name, thisModuleDB)
	exists := false
	_, err := GetUserID(usr.db, name)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		exists = true
	}
	return exists
}

// GetUserMap loads user information from the database and returns it
// inside of a map
func (usr User) GetUserMap(uid int) (map[string]interface{}, error) {
	var userData map[string]interface{}
	var name, email, token string
	err := usr.db.sql.QueryRow(`SELECT name, mail, token FROM users JOIN login_token ON users.uid = login_token.uid AND users.uid = ?`, uid).Scan(&name, &email, &token)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		userData = make(map[string]interface{})
		userData["Name"] = name
		userData["Email"] = email
		userData["Token"] = token
		DebugHandler.Sys("GetUserMap::"+userData["Name"].(string), thisModuleDB)
	}
	return userData, err
}

// CheckToken checks if the user has a token in the database login_token table
func (usr User) CheckToken(uid int) bool {
	var token string
	exists := false
	DebugHandler.Sys("CheckingToken", thisModuleDB)
	err := usr.db.sql.QueryRow(`SELECT token FROM login_token WHERE uid = ?`, uid).Scan(&token)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("TokenNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	default:
		DebugHandler.Sys("TokenFound", thisModuleDB)
		exists = true
	}
	return exists
}

// StoreToken writes user token to the database
func (usr User) StoreToken(uid int, token string) {
	DebugHandler.Sys("StoreToken", thisModuleDB)
	_, err := usr.db.sql.Exec(`INSERT INTO login_token(uid, token) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE token = ?`, uid, token, token)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// GetUserID will return the database uid for a username
func GetUserID(db *Database, name string) (int, error) {
	var uid int
	DebugHandler.Sys("GetUserID::"+name, thisModuleDB)
	err := db.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("UserNotFound", thisModuleDB)
	case err != nil:
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
	return uid, err
}
