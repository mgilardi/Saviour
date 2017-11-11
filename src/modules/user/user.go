package user

import (
      "fmt"
      "modules/database"
      "crypto/rand"
      "encoding/base64"
)

type User struct {
  uid int
  userMap map[string]string
  db *database.Database
  online bool
}

func InitUser(id int, db *database.Database) *User {
  var user User
  var err error
  user.uid = id
  user.db = db
  user.CheckToken()
  user.userMap, err = user.db.GetUserMap(user.uid)
  user.online = true
  if err != nil {
    //
  }
  return &user
}

func (user User) CheckToken() {
  if !user.db.CheckToken(user.uid) {
    user.db.StoreToken(user.uid, genToken(32))
  }
}

func (user User) GetName() string {
  return user.userMap["name"]
}

func (user User) GetToken() string {
  return user.userMap["token"]
}

func (user User) IsOnline() bool {
  return user.online
}
func (user User) LogOn() {
  user.online = true
}
func (user User) LogOff() {
  user.online = false
}

func genToken(length int) string {
  byte := make([]byte, length)
  _, err := rand.Read(byte)
  if err != nil {
    fmt.Println("ErrorGenToken::" + err.Error())
  }
  return base64.URLEncoding.EncodeToString(byte)
}