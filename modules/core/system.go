package core

import (
	//"config"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	thisModuleSystem = "System"
)

// DataPacket is the struct that json files are loaded into when marshaled
type DataPacket struct {
	Login struct {
		User  string `json:"user,omitempty" validate:"max=45,alphanum|len=0"`
		Pass  string `json:"pass,omitempty" validate:"max=45,alphanum|len=0"`
		Email string `json:"email,omitempty" validate:"max=45,email|len=0"`
	} `json:"login"`
	Saviour struct {
		Username string `json:"username,omitempty" validate:"max=45,alphanum|len=0"`
		Status   int    `json:"status,omitempty" validate:"len=3|len=0"`
		Token    string `json:"token,omitempty" validate:"max=45,base64|len=0"`
		Message  string `json:"message,omitempty" validate:"max=45,base64|len=0"`
	} `json:"saviour"`
}

// genDataPacket generates a packet for transaction
func genDataPacket(token string, message string, status int, username string) []byte {
	var packet DataPacket
	var buf []byte
	// Testing
	packet = sanitizePacket(packet)
	packet.Saviour.Token = token
	packet.Saviour.Message = message
	packet.Saviour.Status = status
	packet.Saviour.Username = username
	buf, err := json.Marshal(&packet)
	if err != nil {
		LogHandler.Err(err, thisModuleSystem)
		DebugHandler.Err(err, thisModuleSystem, 3)
	}
	return buf
}

// loadDataPacket loads incoming packet for analysis
func loadDataPacket(buf []byte) DataPacket {
	var packet DataPacket
	err := json.Unmarshal(buf, &packet)
	if err != nil {
		LogHandler.Err(err, thisModuleSystem)
		DebugHandler.Err(err, thisModuleSystem, 3)
	}
	return sanitizePacket(packet)
}

// System contains server responses to http requests
type System struct {
	hostname, port string
	db             *Database
	cache          *Cache
	currentUser    *User
	conUsers       map[string]*User
}

// InitSystem initialize system
func InitSystem(datab *Database) {
	var sys System
	sys.conUsers = make(map[string]*User)
	sys.db = datab
	sys.cache = CacheHandler
	DebugHandler.Sys("Starting", thisModuleSystem)
	exists, options := sys.cache.GetCacheMap("core:config")
	if exists {
		sys.hostname = options["Hostname"].(string)
		sys.port = options["Port"].(string)
		DebugHandler.Sys("LoadedConfigFromCache::"+options["Name"].(string), thisModuleSystem)
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
	DebugHandler.Err(http.ListenAndServe(sys.hostname+":"+sys.port, servRouter), thisModuleSystem, 1)
}

// indexPage handles index page
func (sys *System) indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request_UserLogin:")
	fmt.Println("ReceivedRequest")
}

// createRequest handles user creation/registration
func (sys *System) createRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	loginParam := [3]string{packet.Login.User, packet.Login.Pass, packet.Login.Email}
	DebugHandler.Sys("CreatingUser::"+loginParam[0], thisModuleSystem)
}

// loginRequest handles initial login, if user is not found or password is incorrect it will return a UserNotFound
// error, if user is already in the connected users map and marked as loggedIn the request will return a
// UserAlreadyLoggedIn error, if the username and password is correct it will return the json including the token,
// a status of 200, the username and a message of LoginSuccessful.
func (sys *System) loginRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var userFound, exists bool
	var currentUser *User
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	loginParam := [3]string{packet.Login.User, packet.Login.Pass, packet.Login.Email}
	DebugHandler.Sys("LoginAttempt::"+loginParam[0], thisModuleSystem)
	userFound, currentUser = InitUser(sys.db, loginParam[0], loginParam[1])
	if userFound == true {
		currentUser, exists = sys.conUsers[loginParam[0]]
		if exists && sys.conUsers[loginParam[0]].IsOnline() {
			buf = genDataPacket("", "UserAlreadyLoggedIn", status, loginParam[0])
			DebugHandler.Sys("LoginFailed::UserLoggedIn::"+loginParam[0], thisModuleSystem)
		} else {
			status = 200
			if !exists {
				sys.conUsers[sys.currentUser.GetName()] = currentUser
			}
			sys.currentUser.SetOnline(true)
			buf = genDataPacket(sys.currentUser.GetToken(), "LoginSuccessful", status, sys.currentUser.GetName())
			DebugHandler.Sys("LoginSuccessful::"+sys.currentUser.GetName(), thisModuleSystem)
		}
	} else {
		buf = genDataPacket("", "UserNotFound", status, loginParam[0])
		DebugHandler.Sys("LoginFailed::UserNotFound::"+loginParam[0], thisModuleSystem)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// CreateUser creates a new user entry in the database
func (sys *System) CreateUser(name string, pass string, email string) {
	DebugHandler.Sys("CreateUser::"+name, thisModuleDB)
	_, err := sys.db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	}
}

// RemoveUser removes a user entry from the database
func (sys *System) RemoveUser(name string) {
	DebugHandler.Sys("RemoveUser::"+name, thisModuleDB)
	uid, err := GetUserID(sys.db, name)
	if err != nil {
		LogHandler.Err(err, thisModuleDB)
		DebugHandler.Err(err, thisModuleDB, 3)
	} else {
		tx, err := sys.db.sql.Begin()
		tx.Exec(`DELETE FROM login_token WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM sessions WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM users WHERE uid = ?`, uid)
		if err != nil {
			DebugHandler.Err(err, thisModuleDB, 3)
			tx.Rollback()
		} else {
			tx.Commit()
			sys.db.ResetIncrement("login_token", "user_roles", "sessions", "users")
		}
	}
}
