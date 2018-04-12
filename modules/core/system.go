package core

import (
	//"config"

	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/bcrypt"
)

const (
	// SYSTEM Constant
	SYSTEM = "System"
)

// System contains server responses to http requests
type System struct {
	hostname, port, local string
	conUsers              map[string]*User
}

// InitSystem initialize system
func InitSystem(inProd bool) {
	var sys System
	sys.conUsers = make(map[string]*User)
	Logger("Starting", "System", MSG)
	options := OptionsHandler.GetOption("Core")
	sys.hostname = options["Hostname"].(string)
	sys.port = options["Port"].(string)
	sys.local = options["Location"].(string)
	if !inProd {
		sys.handleRequest()
	} else {
		sys.handleRequestProd()
	}
}

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys System) handleRequest() {
	serv := &http.Server{
		Addr:         sys.hostname + ":" + sys.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	cert, key := sys.loadTLSCert()
	servRouter := mux.NewRouter()
	serv.Handler = servRouter
	servRouter.HandleFunc("/", sys.indexPage)
	servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
	servRouter.HandleFunc("/register", sys.createRequest).Methods("POST")
	servRouter.HandleFunc("/request/logoff", sys.logoffRequest).Methods("POST")
	servRouter.HandleFunc("/request/password", sys.changePassRequest).Methods("POST")

	Logger(serv.ListenAndServeTLS(cert, key).Error(), "System", ERROR)
}

func (sys System) handleRequestProd() {
	certPath, _ := FindPath("cert")
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache(certPath),
		HostPolicy: autocert.HostWhitelist("saviour.diyccs.com"),
	}
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	serv := &http.Server{
		Addr:         sys.hostname + ":" + sys.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	servRouter := mux.NewRouter()
	serv.Handler = servRouter
	servRouter.HandleFunc("/", sys.indexPage)
	servRouter.HandleFunc("/login", sys.loginRequest).Methods("POST")
	servRouter.HandleFunc("/register", sys.createRequest).Methods("POST")
	servRouter.HandleFunc("/request/logoff", sys.logoffRequest).Methods("POST")
	servRouter.HandleFunc("/request/password", sys.changePassRequest).Methods("POST")

	Logger(serv.ListenAndServeTLS("", "").Error(), "System", ERROR)
}

func (sys System) loadTLSCert() (string, string) {
	dir, _ := os.Getwd()
	cert := dir + "/cert/server.crt"
	key := dir + "/cert/server.key"
	return cert, key
}

// indexPage handles index page
func (sys System) indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Request_UserLogin:")
	fmt.Println("ReceivedRequest")
}

// createRequest handles user creation/registration
func (sys System) createRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	status := 400
	switch {
	case !valid:
		Logger("InvalidPacket", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPacket", status, "")
	default:
		status = 200
		Logger("CreatingUser::"+packet.Login.User, SYSTEM, MSG)
		err := CommandHandler.CreateUser(packet.Login.User, packet.Login.Pass, packet.Login.Email)
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
func (sys System) loginRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	Logger("LoginAttempt::"+packet.Login.User, SYSTEM, MSG)
	userFound, userMap, currentUser := InitUser(packet.Login.User, packet.Login.Pass)
	_, exists := sys.conUsers[packet.Login.User]
	switch {
	case !userFound:
		buf = genDataPacket("", "UserNotFound", status, packet.Login.User)
		Logger("LoginFailed::InvalidRequest::"+packet.Login.User, SYSTEM, WARN)
	case !valid:
		buf = genDataPacket("", "InvalidRequest", status, packet.Login.User)
		Logger("LoginFailed::InvalidRequest::"+packet.Login.User, SYSTEM, WARN)
	default:
		status = 200
		if !exists {
			sys.conUsers[packet.Login.User] = currentUser
		}
		buf = genDataPacket(currentUser.GetToken(), "LoginSuccessful", status, userMap["name"].(string))
		Logger("LoginSuccessful::"+userMap["name"].(string), SYSTEM, MSG)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys System) logoffRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !valid:
		Logger("InvalidPackage", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !exists:
		Logger("UserNotConnectedLogoff::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	default:
		userMap := currentUser.GetUserMap()
		if !currentUser.VerifyToken(packet.Saviour.Token) {
			Logger("InvalidTokenLogoff::"+userMap["name"].(string), SYSTEM, WARN)
			buf = genDataPacket("", "InvalidToken", status, userMap["name"].(string))
		} else {
			status = 200
			Logger("LogoffSuccsessful::"+userMap["name"].(string), SYSTEM, MSG)
			buf = genDataPacket(userMap["token"].(string), "LogOff::Sucsessful", status, userMap["name"].(string))
			delete(sys.conUsers, userMap["name"].(string))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys System) changePassRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !valid:
		Logger("InvalidPackage", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !exists:
		Logger("UserNotConnected::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	default:
		userMap := currentUser.GetUserMap()
		if !currentUser.VerifyToken(packet.Saviour.Token) {
			Logger("InvalidTokenChangePassword", SYSTEM, WARN)
			buf = genDataPacket("", "InvalidToken", status, userMap["name"].(string))
		} else {
			status = 200
			Logger("ChangePasswordRequest::"+userMap["name"].(string), SYSTEM, MSG)
			changeRequest := strings.Split(packet.Saviour.Message, ":")
			currentUser.SetPassword(changeRequest[1])
			buf = genDataPacket(userMap["token"].(string), "PasswordChanged", status, userMap["name"].(string))
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys System) changeRoleRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !exists:
		Logger("UserNotConnected::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	case !valid:
		Logger("InvalidPackage", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !currentUser.VerifyToken(packet.Saviour.Token):
		Logger("InvalidToken::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	default:
		changeData := strings.Split(packet.Saviour.Message, ":")
		result := CommandHandler.ChangeUserRole(changeData[0], currentUser, changeData[1])
		buf = genDataPacket(currentUser.GetToken(), result, status, packet.Saviour.Username)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

func (sys System) removeUserRequest(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	status := 400
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	currentUser, exists := sys.conUsers[packet.Saviour.Username]
	switch {
	case !exists:
		Logger("UserNotConnected::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "UserNotConnected", status, packet.Saviour.Username)
	case !valid:
		Logger("InvalidPackage", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	case !currentUser.VerifyToken(packet.Saviour.Token):
		Logger("InvalidToken::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPackage", status, "")
	default:
		result := CommandHandler.RemoveUser(packet.Saviour.Message, currentUser)
		buf = genDataPacket(currentUser.GetToken(), result, status, packet.Saviour.Username)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf)
}

// Admin command contains all adminitstrative commands for the server
var CommandHandler Command

type Command struct {
}

func InitCommand() {
	var cmd Command
	AccessHandler.LoadPerm(cmd)
	CommandHandler = cmd
}

func (cmd Command) PermID() string {
	return "Command"
}

func (cmd Command) DefaultPerm() map[string]bool {
	accessMap := make(map[string]bool)
	return accessMap
}

func (cmd Command) CreateUser(name string, pass string, email string) error {
	var userCheck sql.NullString
	var err error
	Logger("CreateUser::"+name, SYSTEM, MSG)
	DBHandler.sql.QueryRow(`SELECT name FROM users WHERE name = ?`, name).Scan(&userCheck)
	switch {
	case userCheck.Valid:
		Logger("DuplicateUser::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("DuplicateUser::UserCreationFailed")
	case name == "":
		Logger("NameEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("NameEntryIsEmpty::UserCreationFailed")
	case pass == "":
		Logger("PasswordEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("PasswordEntryIsEmpty::UserCreationFailed")
	case email == "":
		Logger("EmailEntryIsEmpty::UserCreationFailed", SYSTEM, WARN)
		err = errors.New("EmailEntryIsEmpty::UserCreationFailed")
	default:
		Logger("CreatingUser::Name:"+name+"::"+email, SYSTEM, MSG)
		hashPass := GenHashPassword(pass)
		roleMap := AccessHandler.genRoleNameMap()
		_, uid := GetUserID(name)
		insertUser := DBHandler.SetupExec(
			`INSERT INTO users (name, pass, mail) `+
				`VALUES (?, ?, ?)`, name, hashPass, email)
		insertUserRole := DBHandler.SetupExec(
			`INSERT INTO user_roles (uid, rid) `+
				`VALUES (?, ?)`, uid, roleMap["user"])
		DBHandler.Exec(insertUser, insertUserRole)
	}
	return err
}

// RemoveUser removes a user entry from the database
func (cmd Command) RemoveUser(removeUser string, requestUser *User) string {
	AccessHandler.AllowAccess("admin", cmd)
	Logger("RemoveUser::"+removeUser, SYSTEM, MSG)
	exists, uid := GetUserID(removeUser)
	switch {
	case !exists:
		Logger("CouldNotRemoveUser::DoesNotExist", SYSTEM, ERROR)
		AccessHandler.DisableAccess("admin", cmd)
		return "UserNotFound"
	case !AccessHandler.CheckPerm(requestUser, cmd):
		Logger("CouldNotRemoverUser::AccessDenied", SYSTEM, WARN)
		AccessHandler.DisableAccess("admin", cmd)
		return "AccessDenied"
	default:
		deleteUserToken := DBHandler.SetupExec(`DELETE FROM login_token WHERE uid = ?`, uid)
		deleteUserRoles := DBHandler.SetupExec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		deleteSessions := DBHandler.SetupExec(`DELETE FROM sessions WHERE uid = ?`, uid)
		deleteUser := DBHandler.SetupExec(`DELETE FROM users WHERE uid = ?`, uid)
		DBHandler.Exec(deleteUserToken, deleteUserRoles, deleteSessions, deleteUser)
		AccessHandler.DisableAccess("admin", cmd)
		return "OperationCompleted::User::" + removeUser + "::Removed"
	}
}

func (cmd Command) ChangeUserRole(changeUser string, requestUser *User, newRole string) string {
	return "OperationCompleted::User::" + changeUser + "::RoleChange::" + newRole
}

// GenHashPassword will hash a password string
func GenHashPassword(pass string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		Logger(err.Error(), SYSTEM, ERROR)
	}
	return string(bytes)
}
