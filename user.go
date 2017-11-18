package main

import (
	"crypto/rand"
	"encoding/base64"
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
func InitUser(id int, db *Database, cache *Cache) *User {
	var user User
	user.uid = id
	user.db = db
	user.cache = cache
	user.CheckToken()
	user.InfoUpdate()
	return &user
}

// CheckToken checks to see if a token exists in the database if not
// it generates one.
func (user *User) CheckToken() {
	if !user.db.CheckToken(user.uid) {
		user.db.StoreToken(user.uid, genToken(32))
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
		userInfo, err = user.db.GetUserMap(user.uid)
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
	token := genToken(32)
	user.db.StoreToken(user.uid, token)
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

// genToken generates user token of specified length
func genToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		LogHandler.Err(err, thisModuleUser)
		DebugHandler.Err(err, thisModuleUser, 1)
	}
	return base64.URLEncoding.EncodeToString(byte)
}
