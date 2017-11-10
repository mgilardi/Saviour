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
        "gopkg.in/go-playground/validator.v9"
        "encoding/json"
        "io/ioutil"
        "regexp"
)

const (
  thisModule = "System"
)

// DataPacket is the struct that json files are loaded into when marshaled
type DataPacket struct {
  Login struct {
    User string `json:"user,omitempty" validate:"max=45,alphanum"`
    Pass string `json:"pass,omitempty" validate:"max=45,alphanum"`
  } `json:"login"`
  Saviour struct {
    Username string `json:"username,omitempty" validate:"max=45,alphanum"`
    Status  string `json:"status,omitempty" validate:"max=3,numeric"`
    Token   string `json:"token,omitempty" validate:"max=45,base64"`
    Message string `json:"message,omitempty" validate:"max=45,base64"`
  } `json:"saviour"`
}

// genDataPacket
func genDataPacket(token string, message string, status string, username string) ([]byte, DataPacket) {
  var packet DataPacket
  var buf []byte
  packet.Saviour.Token = token
  packet.Saviour.Message = message
  packet.Saviour.Status = status
  packet.Saviour.Username = username
  packet = sanitizePacket(packet)
  buf, err := json.Marshal(&packet)
  if err != nil {
    logger.Error("ErrorMarshalPacket::" + err.Error(), thisModule, 3)
  }
  return buf, packet
}

func loadDataPacket(buf []byte) DataPacket {
  var packet DataPacket
  err := json.Unmarshal(buf, &packet)
  if (err != nil) {
    logger.Error("ErrorUnmarshalPacket::" + err.Error(), thisModule, 3)
  }
  return sanitizePacket(packet)
}

func genToken(length int) string {
  byte := make([]byte, length)
  _, err := rand.Read(byte)
  if err != nil {
    fmt.Println("ErrorGenToken::" + err.Error())
  }
  return base64.URLEncoding.EncodeToString(byte)
}

func sanitizeLogin(packet DataPacket) [2]string {
  var login [2]string
  validate := validator.New()
  err := validate.Struct(packet)
  if err != nil {
    regex := regexp.MustCompile("[^a-zA-Z0-9]")
    packet.Login.User = regex.ReplaceAllString(packet.Login.User, "")
    packet.Login.Pass = regex.ReplaceAllString(packet.Login.Pass, "")
    packet.Login.User = checkDataSize(packet.Login.User,45)
    packet.Login.Pass = checkDataSize(packet.Login.Pass,45)
  }
  login[0] = packet.Login.User
  login[1] = packet.Login.Pass
  return login
}

func sanitizePacket(packet DataPacket) DataPacket {
  validate := validator.New()
  err := validate.Struct(packet)
  if err != nil {
    regex := regexp.MustCompile(`\d\d\d{1}`)
    packet.Saviour.Status = regex.FindString(packet.Saviour.Status)
    packet.Saviour.Status = checkDataSize(packet.Saviour.Status,3)
    regex = regexp.MustCompile("[^a-zA-Z0-9]")
    packet.Saviour.Username = regex.ReplaceAllString(packet.Saviour.Username, "")
    packet.Saviour.Username = checkDataSize(packet.Saviour.Username,45)
    regex = regexp.MustCompile(`[^A-Za-z0-9+-\/=]`)
    packet.Saviour.Token = regex.ReplaceAllString(packet.Saviour.Token, "")
    packet.Saviour.Token = checkDataSize(packet.Saviour.Token,45)
    packet.Saviour.Message = regex.ReplaceAllString(packet.Saviour.Message, "")
    packet.Saviour.Message = checkDataSize(packet.Saviour.Message,45)
  }
  return packet
}

func checkDataSize(data string, size int) string {
  var trimStr string
  if (len(data) > size) {
    trimStr = data[0:size]
  } else {
    trimStr = data
  }
  return trimStr
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
  var packet DataPacket
  var loginParam [2]string
  var buf []byte
  buf, _ = ioutil.ReadAll(r.Body)
  packet = loadDataPacket(buf)
  loginParam = sanitizeLogin(packet)
  sys.logger.SystemMessage("LoginAttempt::" + loginParam[0], thisModule)
  if !sys.db.CheckUser(loginParam[0], loginParam[1]) {
    buf, packet = genDataPacket("", "", "401", loginParam[0] )
    sys.logger.Error("LoginFailed", thisModule, 2)
  } else {
    if !sys.db.CheckToken(loginParam[0]) {
      buf, packet = genDataPacket(genToken(32), "Login::Successful", "200", loginParam[0])
      sys.logger.SystemMessage("GeneratedToken::" + packet.Saviour.Token, thisModule)
      sys.db.StoreToken(packet.Saviour.Username, packet.Saviour.Token)
      sys.logger.SystemMessage("LoginSuccessfulGenToken::" + loginParam[0], thisModule)
    } else {
      buf, packet = genDataPacket(sys.db.GetToken(loginParam[0]), "Login::SuccessFul", "200", loginParam[0])
      sys.logger.SystemMessage("LoginSuccessful::" + loginParam[0], thisModule)
    }
  }
  w.Write(buf)
}

