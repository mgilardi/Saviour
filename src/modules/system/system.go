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

type DataPacket struct {
  Login struct {
    User string `json:"user"`
    Pass string `json:"pass"`
  } `json:"login"`
  Saviour struct {
    Username string `json:"username"`
    Status  string `json:"status"`
    Token   string `json:"token"`
    Message string `json:"message"`
  } `json:"saviour"`
}

func genDataPacket(token string, message string, status string, username string) DataPacket {
  var packet DataPacket
  packet.Saviour.Token = token
  packet.Saviour.Message = message
  packet.Saviour.Status = status
  packet.Saviour.Username = username
  return packet
}

func genToken(length int) string {
  byte := make([]byte, length)
  _, err := rand.Read(byte)
  if err != nil {
    fmt.Println("ErrorGenToken::" + err.Error())
  }
  return base64.URLEncoding.EncodeToString(byte)
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
  if (err != nil) {
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
  sys.logger.Error(http.ListenAndServe(":8080", servRouter).Error(), thisModule, 1)
}

func (sys System) indexPage( w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w,"Request_UserLogin:")
  fmt.Println("ReceivedRequest")
}

func (sys System) loginRequest(w http.ResponseWriter, r *http.Request) {
  var err error
  var packet DataPacket
  var buf []byte
  buf, err = ioutil.ReadAll(r.Body)
  err = json.Unmarshal(buf, &packet)
  if (err != nil) {
    fmt.Println(err.Error())
  }
  sys.logger.SystemMessage("LoginAttempt::" + packet.Login.User, thisModule)
  if !sys.db.CheckUser(packet.Login.User, packet.Login.Pass) {
    data := genDataPacket("NULL", "Error::Invalid::Login", "401", packet.Login.User )
    buf, err = json.Marshal(&data)
    sys.logger.Error("LoginFailed", thisModule, 2)
  } else {
    if !sys.db.CheckToken(packet.Login.User) {
      token := genToken(32)
      sys.logger.SystemMessage("GeneratedToken::" + token, thisModule)
      data := genDataPacket(token, "Login::Successful", "200", packet.Login.User)
      sys.db.StoreToken(data.Saviour.Username, data.Saviour.Token)
      buf, err = json.Marshal(&data)
      sys.logger.SystemMessage("LoginSuccessfulGenToken::" + packet.Login.User, thisModule)
    } else {
      data := genDataPacket(sys.db.GetToken(packet.Login.User), "Login::SuccessFul", "200", packet.Login.User)
      buf, err = json.Marshal(&data)
      sys.logger.SystemMessage("LoginSuccessful::" + packet.Login.User, thisModule)
    }
  }
  w.Write(buf)
}

