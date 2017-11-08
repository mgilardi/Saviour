package system

import (
        "config"
        "modules/logger"
        "modules/database"
        "modules/cache"

        "fmt"
        "net/http"
        "crypto/rand"
        "encoding/base64"
        "github.com/gorilla/mux"

        "encoding/json"
        "io/ioutil"
)

const (
  thisModule = "System"
)

type Login struct {
  User string
  Pass string
}

type Token struct {
  Token string
}

func genToken(length int) (error, Token) {
  var t Token
  byte := make([]byte, length)
  _, err := rand.Read(byte)
  t.Token = base64.URLEncoding.EncodeToString(byte)
  return err, t
}

type System struct {
  options *config.Setting
  logger  *logger.LogData
  db *database.Database
}

func InitSystem(conf *[]config.Setting, datab *database.Database, log *logger.LogData, cache *cache.Cache) {
  var sys System
  var err error
  sys.db = datab
  sys.logger = log
  sys.logger.SystemMessage("Starting::Server", thisModule)
  err, sys.options = config.GetSettingModule(thisModule, conf)
  if err != nil {
    //
  }
 sys.startServ()
}

func (sys System) startServ() {
  sys.handleRequest()
}

func (sys System) handleRequest() {
  servRouter := mux.NewRouter()
  servRouter.HandleFunc("/", sys.indexPage)
  servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
  http.ListenAndServe(":8080", servRouter)
}

func (sys System) indexPage( w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w,"Request_UserLogin:")
  fmt.Println("ReceivedRequest")
}

func (sys System) loginRequest( w http.ResponseWriter, r *http.Request) {
  var err error
  var login Login
  var token Token
  var buf []byte
  buf, err = ioutil.ReadAll(r.Body)
  err = json.Unmarshal(buf, &login)
  if err != nil {
    fmt.Println(err.Error())
  }
  if !sys.db.CheckUser(login.User, login.Pass) {
    sys.logger.Error("LoginFailed", thisModule, 2)
  } else {
    sys.logger.SystemMessage("UserVerified::GeneratingToken" + login.User, thisModule)
    err, token = genToken(32)
    if err != nil {
      sys.logger.Error(err.Error(), thisModule, 3)
    }
    sys.logger.SystemMessage("Token::" + token.Token, thisModule)
    loginResponse, err := json.Marshal(&token)
    sys.db.StoreToken(login.User, token.Token)
    if err != nil {
      //
    }
    w.Write(loginResponse)
  }
}

