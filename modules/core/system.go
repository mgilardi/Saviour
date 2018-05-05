package core

import (
	//"config"

	"crypto/tls"
	"database/sql"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
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
	hostname, port, local string
	serv                  *http.Server
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
	sys.handleRequest()

}

// handleRequest sets up router for different webpage requests and redirects them to there function
// ListenAndServ starts the server listing on port
func (sys System) handleRequest() {
	sys.serv = &http.Server{
		Addr:         sys.hostname + ":" + sys.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
		TLSConfig:    &tls.Config{
			//InsecureSkipVerify: true,
		},
	}
	defer sys.serv.Close()
	cert, key := GetCert()
	servRouter := mux.NewRouter()
	sys.serv.Handler = servRouter
	servRouter.PathPrefix("/").HandlerFunc(sys.request)

	Logger(sys.serv.ListenAndServeTLS(cert, key).Error(), "System", ERROR)
}

// Request Function is more general and allows the IPC communication to be checked
// everytime the server responds to a packet being able to interact with the Main
// server thread.
func (sys System) request(w http.ResponseWriter, r *http.Request) {
	var packet DataPacket
	var buf []byte
	var valid bool
	buf, _ = ioutil.ReadAll(r.Body)
	valid, packet = loadDataPacket(buf)
	if !valid {
		Logger("InvalidPacket", SYSTEM, WARN)
		buf = genDataPacket("", "InvalidPacket", 400, "")
	} else {
		switch r.URL.Path {
		case "/login":
			buf = sys.login(packet)
		case "/register":
			buf = sys.create(packet)
		case "/request/password":
			buf = sys.changePass(packet)
		default:
			Logger("InvalidPath", SYSTEM, WARN)
			buf = genDataPacket("", "InvalidPath", 400, "")
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(packet.Saviour.Status)
	w.Write(buf)
	// Handle Agent Signal
	exists, signal := AgentHandler.checkSignal()
	switch {
	case !exists:
		// Do Nothing
	case (signal == "UpdateCert"):
		Logger("UpdateCertAgent::Restarting", SYSTEM, MSG)
		sys.serv.Close()
		os.Exit(0)
	}
}

func (sys System) login(packet DataPacket) []byte {
	var buf []byte
	status := 400
	userFound, userMap, currentUser := InitUser(packet.Login.User, packet.Login.Pass)
	switch {
	case !userFound:
		buf = genDataPacket("", "UserNotFound", status, packet.Login.User)
		Logger("LoginFailed::InvalidRequest::"+packet.Login.User, SYSTEM, WARN)
	default:
		status = 200
		buf = genDataPacket(currentUser.GetToken(), "LoginSuccessful", status, userMap["name"].(string))
		Logger("LoginSuccessful::"+userMap["name"].(string), SYSTEM, MSG)
	}
	return buf
}

func (sys System) create(packet DataPacket) []byte {
	var buf []byte
	status := 400
	Logger("CreatingUser::"+packet.Login.User, SYSTEM, MSG)
	err := CommandHandler.CreateUser(packet.Login.User, packet.Login.Pass, packet.Login.Email)
	if err != nil {
		buf = genDataPacket("", err.Error(), status, packet.Login.User)
	} else {
		status = 200
		buf = genDataPacket("", "UserCreationSucsessful::"+packet.Login.User, status, packet.Login.User)
	}
	return buf
}

func (sys System) changePass(packet DataPacket) []byte {
	var buf []byte
	status := 400
	Logger("ChangingUserPassword::"+packet.Saviour.Username, SYSTEM, MSG)
	exists, userMap, currentUser := GetUser(packet.Saviour.Username, packet.Saviour.Token)
	switch {
	case !exists:
		buf = genDataPacket("", "UserNotAuthenticated", status, packet.Saviour.Username)
	default:
		status = 200
		Logger("ChangePasswordRequest::"+userMap["name"].(string), SYSTEM, MSG)
		changeRequest := strings.Split(packet.Saviour.Message, ":")
		currentUser.SetPassword(changeRequest[1])
		buf = genDataPacket(userMap["token"].(string), "PasswordChanged", status, userMap["name"].(string))
	}
	return buf
}

func (sys System) removeUser(packet DataPacket) []byte {
	var buf []byte
	exists, _, currentUser := GetUser(packet.Saviour.Username, packet.Saviour.Token)
	switch {
	case !exists:
		Logger("UserNotConnected::"+packet.Saviour.Username, SYSTEM, WARN)
		buf = genDataPacket("", "UserNotConnected", 400, packet.Saviour.Username)
	default:
		result := CommandHandler.RemoveUser(packet.Saviour.Message, currentUser)
		buf = genDataPacket(currentUser.GetToken(), result, 200, packet.Saviour.Username)
	}
	return buf
}

// Admin command contains all adminitstrative commands for the server
var CommandHandler Command

type Command struct {
}

func InitCommand() {
	var cmd Command
	CommandHandler = cmd
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
		currentTime := time.Now().Unix()
		insertUser := DBHandler.SetupExec(
			`INSERT INTO users (name, pass, mail, created) `+
				`VALUES (?, ?, ?, ?)`, name, hashPass, email, currentTime)
		DBHandler.Exec(insertUser)
		_, uid := GetUserID(name)
		insertUserRole := DBHandler.SetupExec(
			`INSERT INTO user_roles (uid, rid) `+
				`VALUES (?, ?)`, uid, 3)
		DBHandler.Exec(insertUserRole)
	}
	return err
}

func (cmd Command) changeUserRole(roleChangeUser string, role string, requestUser *User) string {
	Logger("ChangeUserRole::"+roleChangeUser, SYSTEM, MSG)
	//exists, roleChangeUserID := GetUserID(roleChangeUser)
	return "OperationCompleted::User::" + roleChangeUser + "::RoleChange::" + role
}

// RemoveUser removes a user entry from the database
func (cmd Command) RemoveUser(removeUser string, requestUser *User) string {
	Logger("RemoveUser::"+removeUser, SYSTEM, MSG)
	exists, uid := GetUserID(removeUser)
	switch {
	case !exists:
		Logger("CouldNotRemoveUser::DoesNotExist", SYSTEM, ERROR)
		return "UserNotFound"
	default:
		deleteUserToken := DBHandler.SetupExec(`DELETE FROM login_token WHERE uid = ?`, uid)
		deleteUserRoles := DBHandler.SetupExec(`DELETE FROM user_roles WHERE uid = ?`, uid)
		deleteSessions := DBHandler.SetupExec(`DELETE FROM sessions WHERE uid = ?`, uid)
		deleteUser := DBHandler.SetupExec(`DELETE FROM users WHERE uid = ?`, uid)
		DBHandler.Exec(deleteUserToken, deleteUserRoles, deleteSessions, deleteUser)
		return "OperationCompleted::User::" + removeUser + "::Removed"
	}
}

// GenHashPassword will hash a password string
func GenHashPassword(pass string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 14)
	if err != nil {
		Logger(err.Error(), SYSTEM, ERROR)
	}
	return string(bytes)
}

func GetCertServ() (string, string) {
	var newestFile string
	var latestFileTime int64
	var latestCertFile string
	var latestKeyFile string
	options := OptionsHandler.GetOption("Core")
	certFolder := options["CertLocation"].(string) + "csr/"
	keyFolder := options["CertLocation"].(string) + "keys/"
	certfiles, _ := ioutil.ReadDir(certFolder)
	for _, cert := range certfiles {
		file, err := os.Stat(certFolder + cert.Name())
		Logger("FoundCertFile::"+cert.Name(), "System.GetCert", MSG)
		if err != nil {
			Logger("FileStatFailedCert", PACKAGE+"."+"System.GetCert", ERROR)
		}
		if file.ModTime().Unix() > latestFileTime {
			latestFileTime = file.ModTime().Unix()
			newestFile = certFolder + file.Name()
		}
	}
	latestCertFile = newestFile
	latestFileTime = 0
	keyfiles, _ := ioutil.ReadDir(keyFolder)
	for _, key := range keyfiles {
		file, err := os.Stat(keyFolder + key.Name())
		Logger("FoundKeyFile::"+key.Name(), "System.GetCert", MSG)
		if err != nil {
			Logger("FileStatFailedKey", PACKAGE+"."+"System.GetCert", ERROR)
		}
		if file.ModTime().Unix() > latestFileTime {
			latestFileTime = file.ModTime().Unix()
			newestFile = keyFolder + file.Name()
		}
	}
	latestKeyFile = newestFile
	Logger("LoadedCert::"+latestCertFile, SYSTEM, MSG)
	Logger("LoadedKey::"+latestKeyFile, SYSTEM, MSG)
	return latestCertFile, latestKeyFile
}

func GetCert() (string, string) {
	options := OptionsHandler.GetOption("Core")
	certFolder := options["CertLocation"].(string)
	return certFolder + "/fullchain.pem", certFolder + "/privkey.pem"
}
