package system

import (
        "config"
        "modules/logger"
        "modules/database"
        "modules/cache"
        "modules/user"
        "fmt"
        "net/http"
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
    Email string `json:"email,omitempty" validate:"max=45,email"`
  } `json:"login"`
  Saviour struct {
    Username string `json:"username,omitempty" validate:"max=45,alphanum"`
    Status  int `json:"status,omitempty" validate:"max=3"`
    Token   string `json:"token,omitempty" validate:"max=45,base64"`
    Message string `json:"message,omitempty" validate:"max=45,base64"`
  } `json:"saviour"`
}

// genDataPacket generates a packet for transaction
func genDataPacket(token string, message string, status int, username string) []byte {
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
  return buf
}

// loadDataPacket loads incoming packet for analysis
func loadDataPacket(buf []byte) DataPacket {
  var packet DataPacket
  err := json.Unmarshal(buf, &packet)
  if (err != nil) {
    logger.Error("ErrorUnmarshalPacket::" + err.Error(), thisModule, 3)
  }
  return sanitizePacket(packet)
}

// sanitizeLogin validates login packet if fail sanitize
func sanitizeLogin(packet DataPacket) [3]string {
  var login [3]string
  validate := validator.New()
  err := validate.Struct(packet)
  if err != nil {
    regex := regexp.MustCompile("[^a-zA-Z0-9]")
    packet.Login.User = regex.ReplaceAllString(packet.Login.User, "")
    packet.Login.Pass = regex.ReplaceAllString(packet.Login.Pass, "")
    regex = regexp.MustCompile("[^a-zA-Z0-9._@]")
    packet.Login.Email = regex.ReplaceAllString(packet.Login.Email, "")
    packet.Login.User = checkDataSize(packet.Login.User,45)
    packet.Login.Pass = checkDataSize(packet.Login.Pass,45)
    packet.Login.Email = checkDataSize(packet.Login.Email,45)
  }
  login[0] = packet.Login.User
  login[1] = packet.Login.Pass
  login[2] = packet.Login.Email
  return login
}

// sanitizePacket validates regular data packet if fail sanitize
func sanitizePacket(packet DataPacket) DataPacket {
  validate := validator.New()
  err := validate.Struct(packet)
  if err != nil {
    regex := regexp.MustCompile("[^a-zA-Z0-9]")
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

// checkDataSize trims down strings if they are above validation size
func checkDataSize(data string, size int) string {
  var trimStr string
  if (len(data) > size) {
    trimStr = data[0:size]
  } else {
    trimStr = data
  }
  return trimStr
}


// System contains server responses to http requests
type System struct {
  options *config.Setting
  logger  *logger.LogData
  db *database.Database
  conUsers map[string]*user.User
}

// InitSystem initialize system
func InitSystem(conf *[]config.Setting, datab *database.Database, log *logger.LogData, cache *cache.Cache) {
  var sys System
  var err error
  sys.conUsers = make(map[string]*user.User)
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

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys System) handleRequest() {
  servRouter := mux.NewRouter()
  servRouter.HandleFunc("/", sys.indexPage)
  servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
  sys.logger.Error(http.ListenAndServe(":8080", servRouter).Error(), thisModule, 1)
}

// indexPage handles index page
func (sys System) indexPage( w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w,"Request_UserLogin:")
  fmt.Println("ReceivedRequest")
}

// createRequest handles user creation/registration
func (sys System) createRequest(w http.ResponseWriter, r *http.Request) {
  var packet DataPacket
  var loginParam [3]string
  var buf []byte
  buf, _ = ioutil.ReadAll(r.Body)
  packet = loadDataPacket(buf)
  loginParam = sanitizeLogin(packet)
  sys.logger.SystemMessage("CreatingUser::" + loginParam[0], thisModule)
}

// loginRequest handles initial login, if user is not found or password is incorrect it will return a UserNotFound
// error, if user is already in the connected users map and marked as loggedIn the request will return a
// UserAlreadyLoggedIn error, if the username and password is correct it will return the json including the token,
// a status of 100, the username and a message of LoginSuccessful.
func (sys System) loginRequest(w http.ResponseWriter, r *http.Request) {
  var packet DataPacket
  var loginParam [3]string
  var buf []byte
  status := 400
  buf, _ = ioutil.ReadAll(r.Body)
  packet = loadDataPacket(buf)
  loginParam = sanitizeLogin(packet)
  sys.logger.SystemMessage("LoginAttempt::" + loginParam[0], thisModule)
  userFound, uid := sys.db.CheckUserLogin(loginParam[0], loginParam[1])
  if userFound == true {
    currentUser, exists := sys.conUsers[loginParam[0]]
    if exists && sys.conUsers[loginParam[0]].IsOnline() {
      buf = genDataPacket("", "UserAlreadyLoggedIn", status, loginParam[0])
      sys.logger.SystemMessage("LoginFailed::UserLoggedIn::"+loginParam[0], thisModule)
    } else {
      status = 100
      if !exists {
        currentUser = user.InitUser(uid, sys.db)
        sys.conUsers[currentUser.GetName()] = currentUser
      }
      currentUser.LogOn()
      buf = genDataPacket(currentUser.GetToken(), "LoginSuccessful", status, currentUser.GetName())
      sys.logger.SystemMessage("LoginSuccessful::"+currentUser.GetName(), thisModule)
    }
  } else {
    buf = genDataPacket("", "UserNotFound", status, loginParam[0])
    sys.logger.SystemMessage("LoginFailed::UserNotFound::" + loginParam[0], thisModule)
  }
  w.WriteHeader(status)
  w.Write(buf)
}

