package core

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User handles users
type User struct {
	uid         int
	name, token string
	db          *Database
	online      bool
}

// InitUser constructs the user on initial login
// CheckUserLogin inputs username and password and checks if it matches
// the database returns a boolean if the user if found along with the
// users database uid
func InitUser(name string, pass string) (bool, map[string]interface{}, *User) {
	var verified = false
	var user User
	var exists bool
	var cacheMap map[string]interface{}
	user.online = false
	user.db = DBHandler
	Sys("CheckUserLogin::"+name, "User")
	exists, user.uid = GetUserID(name)
	if exists {
		cacheMap = user.GetCache()
		if user.CheckPassword(pass, cacheMap["pass"].(string)) {
			user.token = user.CheckTokenExists()
			verified = true
		} else {
			Warn(errors.New("InvalidPassword::"+name), "User")
		}
	}
	return verified, cacheMap, &user
}

// CheckTokenExists checks to see if a token exists in the database if not
// it generates one.
func (user *User) CheckTokenExists() string {
	exists, token := user.CheckToken(user.uid)
	if !exists {
		token = GenToken(32)
		user.StoreToken(user.uid, token)
		user.UpdateCache()
	}
	return token
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
	return user.token
}

// IsOnline returns the online flag for the user
func (user *User) IsOnline() bool {
	return user.online
}

// GetUserMap returns the user map containing the database user information
func (user *User) GetUserMap() map[string]interface{} {
	return user.GetCache()
}

// SetOnline will set the flag to the input
func (user *User) SetOnline(isOnline bool) {
	if !isOnline {
		Sys("SetOffline", "User")
		user.online = false
		offline := "Offline"
		_, err := user.db.sql.Exec(`UPDATE users SET status = ? WHERE uid = ?`, offline, user.uid)
		if err != nil {
			Error(err, "User")
		}
	} else {
		Sys("SetOnline", "User")
		online := "Online"
		user.online = true
		_, err := user.db.sql.Exec(`UPDATE users SET status = ? WHERE uid = ?`, online, user.uid)
		if err != nil {
			Error(err, "User")
		}
	}
}

// SetPassword sets the users password
func (user *User) SetPassword(pass string) {
	hashPass := GenHashPassword(pass)
	_, err := user.db.sql.Exec(`UPDATE users SET pass = ? WHERE uid = ?`, hashPass, user.uid)
	if err != nil {
		Error(err, "User")
	}
	user.UpdateCache()
}

// CheckPassword checks input password with the hash stored in the database
func (user *User) CheckPassword(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return false
	}
	return true
}

// CheckToken checks if the user has a token in the database login_token table
func (user *User) CheckToken(uid int) (bool, string) {
	var token string
	var expTime int64
	exists := false
	err := user.db.sql.QueryRow(`SELECT token, expires FROM login_token WHERE uid = ?`, uid).Scan(&token, &expTime)
	switch {
	case err == sql.ErrNoRows:
		Sys("TokenNotFound", "User")
	case err != nil:
		Error(err, "User")
	default:
		Sys("TokenFound", "User")
		if time.Now().Unix() < expTime {
			exists = true
		}
	}
	return exists, token
}

// StoreToken writes user token to the database
func (user *User) StoreToken(uid int, token string) {
	Sys("StoreNewToken", "User")
	expTime := time.Now().Add(time.Duration(24) * time.Hour).Unix()
	_, err := user.db.sql.Exec(`INSERT INTO login_token(uid, token, expires) VALUES (?, ?, ?)`+
		`ON DUPLICATE KEY UPDATE token = ?, expires = ?`, uid, token, expTime, token, expTime)
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
	var uid, access int
	var name, pass, email, status, token string
	cacheMap := make(map[string]interface{})
	err := user.db.sql.QueryRow(
		`SELECT users.uid, users.name, users.pass, users.mail, users.status, login_token.token, role.weight FROM users `+
			`INNER JOIN login_token ON users.uid = login_token.uid `+
			`INNER JOIN user_roles ON user_roles.uid = users.uid `+
			`INNER JOIN role ON user_roles.rid = role.rid `+
			`WHERE users.uid = ?`, user.uid).Scan(&uid, &name, &pass, &email, &status, &token, &access)
	if err != nil {
		Error(err, "User")
	}
	cacheMap["uid"] = uid
	cacheMap["name"] = name
	cacheMap["pass"] = pass
	cacheMap["email"] = email
	cacheMap["token"] = token
	cacheMap["status"] = status
	cacheMap["level"] = access
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
func GetUserID(name string) (bool, int) {
	var uid int
	exists := true
	Sys("GetUserID::"+name, "User")
	err := DBHandler.sql.QueryRow("SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		Sys("UserNotFound", "User")
		exists = false
	case err != nil:
		Error(err, "User")
		exists = false
	}
	return exists, uid
}
