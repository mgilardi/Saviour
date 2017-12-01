package core

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User handles users
type User struct {
	uid         int
	roles       map[string]int
	name, token string
	online      bool
}

const (
	// USER module name constant
	USER = "User"
)

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
	Logger("CheckUserLogin::"+name, USER, MSG)
	exists, user.uid = GetUserID(name)
	if exists {
		user.UpdateCache()
		cacheMap = user.GetCache()
		if user.CheckPassword(pass, cacheMap["pass"].(string)) {
			user.token = user.CheckTokenExists()
			user.name = cacheMap["name"].(string)
			user.roles = user.GetUserRoleMap()
			verified = true
		} else {
			Logger("InvalidPassword::"+name, USER, WARN)
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
	exists := false
	userToken := user.GetToken()
	Logger("Match::Token::"+token+"::"+userToken, USER, MSG)
	if token == userToken {
		exists = true
	}
	return exists
}

// GetName returns username
func (user *User) GetName() string {
	return user.name
}

// GetRoleNames returns the role names assigned to this user
func (user *User) GetRoleNames() map[string]int {
	return user.roles
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

// GetUserRoleMap is loading the role map from the database and returning it
func (user *User) GetUserRoleMap() map[string]int {
	var rid int
	var roleName string
	var roleMap map[string]int
	rows, err := DBHandler.sql.Query(`SELECT user_roles.rid, role.name FROM users `+
		`INNER JOIN user_roles ON users.uid = user_roles.uid `+
		`INNER JOIN role ON user_roles.rid = role.rid `+
		`WHERE users.uid = ?`, user.uid)
	if err != nil {
		Logger(err.Error(), USER, ERROR)
	}
	roleMap = make(map[string]int)
	for rows.Next() {
		rows.Scan(&rid, &roleName)
		roleMap[roleName] = rid
	}
	return roleMap
}

// SetOnline will set the flag to the input
func (user *User) SetOnline(isOnline bool) {
	if !isOnline {
		Logger("SetOffline", USER, MSG)
		user.online = false
		offline := "Offline"
		_, err := DBHandler.sql.Exec(
			`UPDATE users SET status = ? `+
				`WHERE uid = ?`, offline, user.uid)
		if err != nil {
			Logger(err.Error(), USER, ERROR)
		}
	} else {
		Logger("SetOnline", USER, MSG)
		online := "Online"
		user.online = true
		_, err := DBHandler.sql.Exec(
			`UPDATE users SET status = ? `+
				`WHERE uid = ?`, online, user.uid)
		if err != nil {
			Logger(err.Error(), USER, ERROR)
		}
	}
}

// SetPassword sets the users password
func (user *User) SetPassword(pass string) {
	hashPass := GenHashPassword(pass)
	_, err := DBHandler.sql.Exec(
		`UPDATE users SET pass = ? `+
			`WHERE uid = ?`, hashPass, user.uid)
	if err != nil {
		Logger(err.Error(), USER, ERROR)
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
	err := DBHandler.sql.QueryRow(
		`SELECT token, expires FROM login_token `+
			`WHERE uid = ?`, uid).Scan(&token, &expTime)
	switch {
	case err == sql.ErrNoRows:
		Logger("TokenNotFound", USER, MSG)
	case err != nil:
		Logger(err.Error(), USER, ERROR)
	default:
		if time.Now().Unix() < expTime {
			Logger("TokenFound", USER, MSG)
			exists = true
		}
	}
	return exists, token
}

// StoreToken writes user token to the database
func (user *User) StoreToken(uid int, token string) {
	Logger("StoreNewToken", USER, MSG)
	expTime := time.Now().Add(time.Duration(24) * time.Hour).Unix()
	_, err := DBHandler.sql.Exec(
		`INSERT INTO login_token(uid, token, expires) VALUES (?, ?, ?) `+
			`ON DUPLICATE KEY UPDATE token = ?, expires = ?`, uid, token, expTime, token, expTime)
	if err != nil {
		Logger(err.Error(), USER, ERROR)
	}
}

// GetCache requests user info from cache
func (user *User) GetCache() map[string]interface{} {
	exists, userCache := CacheHandler.Cache(user)
	if !exists {
		Logger("FailedToLoadCache", USER, ERROR)
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
	var name, pass, email, status, token string
	cacheMap := make(map[string]interface{})
	err := DBHandler.sql.QueryRow(
		`SELECT users.uid, users.name, users.pass, users.mail, users.status FROM users `+
			`WHERE users.uid = ?`, user.uid).Scan(&uid, &name, &pass, &email, &status)
	if err != nil {
		Logger(err.Error(), USER, ERROR)
	}
	Logger("UserCache::"+name+"::"+pass+"::"+email+"::"+status+"::"+token, USER, MSG)
	cacheMap["uid"] = uid
	cacheMap["name"] = name
	cacheMap["pass"] = pass
	cacheMap["email"] = email
	cacheMap["status"] = status
	cacheMap["token"] = token
	return cacheMap
}

// GenToken generates user token of specified length
func GenToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		Logger(err.Error(), USER, ERROR)
	}
	return base64.URLEncoding.EncodeToString(byte)
}

// GetUserID will return the database uid for a username
func GetUserID(name string) (bool, int) {
	var uid int
	exists := true
	Logger("GetUserID::"+name, USER, MSG)
	err := DBHandler.sql.QueryRow(
		"SELECT uid FROM users WHERE name = ?", name).Scan(&uid)
	switch {
	case err == sql.ErrNoRows:
		Logger("UserNotFound", USER, MSG)
		exists = false
	case err != nil:
		Logger(err.Error(), USER, ERROR)
		exists = false
	}
	return exists, uid
}
