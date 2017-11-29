package core

import (
	//"config"

	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SYSTEM Constant
	SYSTEM = "System"
)

// System contains server responses to http requests
type System struct {
	hostname, port string
	db             *Database
	conUsers       map[string]*User
}

// InitSystem initialize system
func InitSystem() {
	var sys System
	sys.conUsers = make(map[string]*User)
	sys.db = DBHandler
	Sys("Starting", SYSTEM)
	options := OptionsHandler.GetOptions("core")
	sys.hostname = options["Hostname"].(string)
	sys.port = options["Port"].(string)
	sys.handleRequest()
}

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys *System) handleRequest() {
	serv := &http.Server{
		Addr:         sys.hostname + ":" + sys.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	servRouter := mux.NewRouter()
	serv.Handler = servRouter
	servRouter.HandleFunc("/", sys.indexPage)
	servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
	servRouter.HandleFunc("/register", sys.createRequest).Methods("POST")
	servRouter.HandleFunc("/request/logoff", sys.logoffRequest).Methods("POST")
	servRouter.HandleFunc("/request/password", sys.changePassRequest).Methods("POST")

	Error(serv.ListenAndServe(), SYSTEM)
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
	var valid bool
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	status := 400
	switch {
	case !valid:
		Error(errors.New("InvalidPacket"), SYSTEM)
		buf = genDataPacket("", "InvalidPacket", status, "")
	default:
		status = 200
		Sys("CreatingUser::"+packet.Login.User, SYSTEM)
		err := sys.CreateUser(packet.Login.User, packet.Login.Pass, packet.Login.Email)
		if err != nil {
			status = 400
			buf = genDataPacket("", err.Error(), status, packet.Login.User)
		} else {
			buf = genDataPacket("", "UserCreationSucsessful::"+packet.Login.User, status, packet.Login.User)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// loginRequest handles initial login, if user is not found or password is incorrect it will return a UserNotFound
// error, if user is already in the connected users map and marked as loggedIn the request will return a
// UserAlreadyLoggedIn error, if the username and password is correct it will return the json including the token,
// a status of 200, the username and a message of LoginSuccessful.
func (sys *System) loginRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	Sys("LoginAttempt::"+packet.Login.User, SYSTEM)
	userFound, userMap, currentUser := InitUser(packet.Login.User, packet.Login.Pass)
	_, exists := sys.conUsers[packet.Login.User]
	switch {
	case !userFound:
		buf = genDataPacket("", "UserNotFound", status, packet.Login.User)
		Sys("LoginFailed::InvalidRequest::"+packet.Login.User, SYSTEM)
	case !valid:
		buf = genDataPacket("", "InvalidRequest", status, packet.Login.User)
		Sys("LoginFailed::InvalidRequest::"+packet.Login.User, SYSTEM)
	case exists:
		buf = genDataPacket("", "UserAlreadyLoggedIn", status, userMap["name"].(string))
		Sys("LoginFailed::UserLoggedIn::"+userMap["name"].(string), SYSTEM)
	default:
		status = 200
		currentUser.SetOnline(true)
		sys.conUsers[packet.Login.User] = currentUser
		buf = genDataPacket(userMap["token"].(string), "LoginSuccessful", status, userMap["name"].(string))
		Sys("LoginSuccessful::"+userMap["name"].(string), SYSTEM)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys *System) logoffRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !valid:
		Sys("InvalidPackage", SYSTEM)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !exists:
		Sys("UserNotConnectedLogoff::"+packet.Saviour.Username, SYSTEM)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	default:
		userMap := currentUser.GetUserMap()
		if userMap["token"].(string) != packet.Saviour.Token {
			Sys("InvalidTokenLogoff::"+userMap["name"].(string), SYSTEM)
			buf = genDataPacket("", "InvalidToken", status, userMap["name"].(string))
		} else {
			status = 200
			currentUser.SetOnline(false)
			Sys("LogoffSuccsessful::"+userMap["name"].(string), SYSTEM)
			buf = genDataPacket(userMap["token"].(string), "LogOff::Sucsessful", status, userMap["name"].(string))
			delete(sys.conUsers, userMap["name"].(string))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys *System) changePassRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !valid:
		Sys("InvalidPackage", SYSTEM)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !exists:
		Sys("UserNotConnected::"+packet.Saviour.Username, SYSTEM)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	default:
		userMap := currentUser.GetUserMap()
		if userMap["token"].(string) != packet.Saviour.Token {
			Sys("InvalidTokenChangePassword", SYSTEM)
			buf = genDataPacket("", "InvalidToken", status, userMap["name"].(string))
		} else {
			status = 200
			Sys("ChangePasswordRequest::"+userMap["name"].(string), SYSTEM)
			changeRequest := strings.Split(packet.Saviour.Message, ":")
			currentUser.SetPassword(changeRequest[1])
			buf = genDataPacket(userMap["token"].(string), "PasswordChanged", status, userMap["name"].(string))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// CreateUser creates a new user entry in the database
func (sys *System) CreateUser(name string, pass string, email string) error {
	var userCheck sql.NullString
	var err error
	Sys("CreateUser::"+name, SYSTEM)
	sys.db.sql.QueryRow(`SELECT name FROM users WHERE name = ?`, name).Scan(&userCheck)
	switch {
	case userCheck.Valid:
		Sys("DuplicateUser::UserCreationFailed", SYSTEM)
		err = errors.New("DuplicateUser::UserCreationFailed")
	case name == "":
		Sys("NameEntryIsEmpty::UserCreationFailed", SYSTEM)
		err = errors.New("NameEntryIsEmpty::UserCreationFailed")
	case pass == "":
		Sys("PasswordEntryIsEmpty::UserCreationFailed", SYSTEM)
		err = errors.New("PasswordEntryIsEmpty::UserCreationFailed")
	case email == "":
		Sys("EmailEntryIsEmpty::UserCreationFailed", SYSTEM)
		err = errors.New("EmailEntryIsEmpty::UserCreationFailed")
	default:
		hashPass := GenHashPassword(pass)
		_, dberr := sys.db.sql.Exec(`INSERT INTO users (name, pass, mail) VALUES (?, ?, ?)`, name, hashPass, email)
		if err != nil {
			Error(dberr, SYSTEM)
		}
	}
	return err
}

// RemoveUser removes a user entry from the database
func (sys *System) RemoveUser(name string) {
	Sys("RemoveUser::"+name, SYSTEM)
	exists, uid := GetUserID(name)
	if exists {
		tx, err := sys.db.sql.Begin()
		tx.Exec(`DELETE FROM login_token WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM sessions WHERE uid = ?`, uid)
		tx.Exec(`DELETE FROM users WHERE uid = ?`, uid)
		if err != nil {
			Error(err, SYSTEM)
			tx.Rollback()
		} else {
			tx.Commit()
			sys.db.ResetIncrement("login_token", "user_roles", "sessions", "users")
		}
	} else {
		Error(errors.New("CouldNotRemoveUser::DoesNotExist"), SYSTEM)
	}
}

// GenHashPassword will hash a password string
func GenHashPassword(pass string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		Error(err, "User")
	}
	return string(bytes)
}
