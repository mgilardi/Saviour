package user

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"modules/database"
	"strconv"
)

type User struct {
	uid                int
	name, token, email string
	db                 *database.Database
	cache              *database.Cache
	online             bool
}

func InitUser(id int, db *database.Database, cache *database.Cache) *User {
	var user User
	user.uid = id
	user.db = db
	user.cache = cache
	user.CheckToken()
	user.InfoUpdate()
	return &user
}

func (user *User) CheckToken() {
	if !user.db.CheckToken(user.uid) {
		user.db.StoreToken(user.uid, genToken(32))
	}
}

func (user *User) InfoUpdate() {
	userInfo := user.GetInfoMap()
	user.name = userInfo["Name"].(string)
	user.token = userInfo["Token"].(string)
	user.email = userInfo["Email"].(string)
}

func (user *User) GetInfoMap() map[string]interface{} {
	exists, userInfo := user.cache.GetCacheMap("user:" + strconv.Itoa(user.uid) + ":" + "info")
	if !exists {
		var err error
		userInfo, err = user.db.GetUserMap(user.uid)
		if err != nil {

		}
		user.cache.SetCacheMap("user:"+strconv.Itoa(user.uid)+":"+"info", userInfo, false)
	}
	return userInfo
}

func (user *User) GetName() string {
	return user.name
}

func (user *User) GetToken() string {
	return user.token
}

func (user *User) GetEmail() string {
	return user.email
}

func (user *User) SetToken() {
	token := genToken(32)
	user.db.StoreToken(user.uid, token)
}

func (user *User) IsOnline() bool {
	return user.online
}

func (user *User) SetOnline(isOnline bool) {
	user.online = isOnline
}

func genToken(length int) string {
	byte := make([]byte, length)
	_, err := rand.Read(byte)
	if err != nil {
		fmt.Println("ErrorGenToken::" + err.Error())
	}
	return base64.URLEncoding.EncodeToString(byte)
}
