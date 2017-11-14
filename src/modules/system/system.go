package system

import (
	//"config"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"modules/database"
	"modules/logger"
	"modules/user"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	thisModule = "System"
)

// DataPacket is the struct that json files are loaded into when marshaled
type DataPacket struct {
	Login struct {
		User  string `json:"user,omitempty" validate:"min=0,max=45,alphanum"`
		Pass  string `json:"pass,omitempty" validate:"min=0,max=45,alphanum"`
		Email string `json:"email,omitempty" validate:"min=0,max=45,email"`
	} `json:"login"`
	Saviour struct {
		Username string `json:"username,omitempty" validate:"max=45,alphanum"`
		Status   int    `json:"status,omitempty" validate:"max=3"`
		Token    string `json:"token,omitempty" validate:"max=45,base64"`
		Message  string `json:"message,omitempty" validate:"max=45,base64"`
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
		logger.Error("ErrorMarshalPacket::"+err.Error(), thisModule, 3)
	}
	return buf
}

// loadDataPacket loads incoming packet for analysis
func loadDataPacket(buf []byte) DataPacket {
	var packet DataPacket
	err := json.Unmarshal(buf, &packet)
	if err != nil {
		logger.Error("ErrorUnmarshalPacket::"+err.Error(), thisModule, 3)
	}
	return sanitizePacket(packet)
}

// System contains server responses to http requests
type System struct {
	hostname, port string
	db             *database.Database
	cache          *database.Cache
	currentUser    *user.User
	conUsers       map[string]*user.User
}

// InitSystem initialize system
func InitSystem(datab *database.Database, cache *database.Cache) {
	var sys System
	sys.conUsers = make(map[string]*user.User)
	sys.db = datab
	sys.cache = cache
	logger.SystemMessage("Starting::Server", thisModule)
	exists, options := sys.cache.GetCacheMap("system:config")
	if exists {
		sys.hostname = options["Hostname"].(string)
		sys.port = options["Port"].(string)
		logger.SystemMessage("LoadedConfigFromCache::"+options["Name"].(string), thisModule)
	}
	sys.startServ()
}

func (sys *System) startServ() {
	sys.handleRequest()
}

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys *System) handleRequest() {
	servRouter := mux.NewRouter()
	servRouter.HandleFunc("/", sys.indexPage)
	servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
	logger.Error(http.ListenAndServe(sys.hostname+":"+sys.port, servRouter).Error(), thisModule, 1)
}

// indexPage handles index page
func (sys *System) indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request_UserLogin:")
	fmt.Println("ReceivedRequest")
}

// createRequest handles user creation/registration
func (sys *System) createRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var loginParam [3]string
	var buf []byte
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	loginParam = sanitizeLogin(packet)
	logger.SystemMessage("CreatingUser::"+loginParam[0], thisModule)
}

// loginRequest handles initial login, if user is not found or password is incorrect it will return a UserNotFound
// error, if user is already in the connected users map and marked as loggedIn the request will return a
// UserAlreadyLoggedIn error, if the username and password is correct it will return the json including the token,
// a status of 100, the username and a message of LoginSuccessful.
func (sys *System) loginRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var loginParam [3]string
	var buf []byte
	var uid int
	var userFound, exists bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	loginParam = sanitizeLogin(packet)
	logger.SystemMessage("LoginAttempt::"+loginParam[0], thisModule)
	userFound, uid = sys.db.CheckUserLogin(loginParam[0], loginParam[1])
	if userFound == true {
		sys.currentUser, exists = sys.conUsers[loginParam[0]]
		if exists && sys.conUsers[loginParam[0]].IsOnline() {
			buf = genDataPacket("", "UserAlreadyLoggedIn", status, loginParam[0])
			logger.SystemMessage("LoginFailed::UserLoggedIn::"+loginParam[0], thisModule)
		} else {
			status = 200
			if !exists {
				sys.currentUser = user.InitUser(uid, sys.db, sys.cache)
				sys.conUsers[sys.currentUser.GetName()] = sys.currentUser
			}
			sys.currentUser.SetOnline(true)
			buf = genDataPacket(sys.currentUser.GetToken(), "LoginSuccessful", status, sys.currentUser.GetName())
			logger.SystemMessage("LoginSuccessful::"+sys.currentUser.GetName(), thisModule)
		}
	} else {
		buf = genDataPacket("", "UserNotFound", status, loginParam[0])
		logger.SystemMessage("LoginFailed::UserNotFound::"+loginParam[0], thisModule)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}
