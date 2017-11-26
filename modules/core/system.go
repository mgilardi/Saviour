package core

import (
	//"config"

	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

// System contains server responses to http requests
type System struct {
	hostname, port string
	db             *Database
	cache          *Cache
	conUsers       map[string]*User
}

// InitSystem initialize system
func InitSystem() {
	var sys System
	sys.conUsers = make(map[string]*User)
	sys.db = DBHandler
	Sys("Starting", "System")
	options := OptionsHandler.GetOptions("core")
	sys.hostname = options["Hostname"].(string)
	sys.port = options["Port"].(string)
	sys.handleRequest()
}

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys *System) handleRequest() {
	servRouter := mux.NewRouter()
	servRouter.HandleFunc("/", sys.indexPage)
	servRouter.HandleFunc("/logoff", sys.logoffRequest).Methods("POST")
	servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
	Error(http.ListenAndServe(sys.hostname+":"+sys.port, servRouter), "System")
}

// indexPage handles index page
func (sys *System) indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request_UserLogin:")
	fmt.Println("ReceivedRequest")
}

func (sys *System) logoffRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	if !exists {
		Sys("UserNotConnectedLogoff::"+packet.Saviour.Username, "System")
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	} else {
		if !currentUser.VerifyToken(packet.Saviour.Token) {
			Sys("InvalidTokenLogoff::"+currentUser.GetName(), "System")
			buf = genDataPacket("", "InvalidToken", status, packet.Saviour.Username)
		} else {
			status = 200
			Sys("LogoffSuccsessful::"+currentUser.GetName(), "System")
			buf = genDataPacket(currentUser.GetToken(), "LogOff::Sucsessful", status, currentUser.GetName())
			delete(sys.conUsers, currentUser.GetName())
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// createRequest handles user creation/registration
func (sys *System) createRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	Sys("CreatingUser::"+packet.Login.User, "System")
}

// loginRequest handles initial login, if user is not found or password is incorrect it will return a UserNotFound
// error, if user is already in the connected users map and marked as loggedIn the request will return a
// UserAlreadyLoggedIn error, if the username and password is correct it will return the json including the token,
// a status of 200, the username and a message of LoginSuccessful.
func (sys *System) loginRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var exists bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	packet = loadDataPacket(buf)
	Sys("LoginAttempt::"+packet.Login.User, "System")
	userFound, currentUser := InitUser(packet.Login.User, packet.Login.Pass)
	if userFound == true {
		_, exists = sys.conUsers[packet.Login.User]
		if exists && sys.conUsers[currentUser.GetName()].IsOnline() {
			buf = genDataPacket("", "UserAlreadyLoggedIn", status, packet.Login.User)
			Sys("LoginFailed::UserLoggedIn::"+packet.Login.User, "System")
		} else {
			status = 200
			currentUser.SetOnline(true)
			if !exists {
				sys.conUsers[currentUser.GetName()] = currentUser
			}
			buf = genDataPacket(currentUser.GetToken(), "LoginSuccessful", status, currentUser.GetName())
			Sys("LoginSuccessful::"+currentUser.GetName(), "System")
		}
	} else {
		buf = genDataPacket("", "UserNotFound", status, packet.Login.User)
		Sys("LoginFailed::UserNotFound::"+packet.Login.User, "System")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// CreateUser creates a new user entry in the database
func (sys *System) CreateUser(name string, pass string, email string) {
	Sys("CreateUser::"+name, "System")
	_, err := sys.db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, pass, email)
	if err != nil {
		Error(err, "System")
	}
}

// RemoveUser removes a user entry from the database
func (sys *System) RemoveUser(name string) {
	Sys("RemoveUser::"+name, "System")
	uid, err := GetUserID(sys.db, name)
	if err != nil {
		Error(err, "System")
	} else {
		tx, err := sys.db.sql.Begin()
		tx.Exec(`DELETE FROM login_token WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM sessions WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM users WHERE uid = ?`, uid)
		if err != nil {
			Error(err, "System")
			tx.Rollback()
		} else {
			tx.Commit()
			sys.db.ResetIncrement("login_token", "user_roles", "sessions", "users")
		}
	}
}
