package core

import (
	//"config"

	"crypto/tls"
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
		case "request/removeuser":
			buf = sys.removeUser(packet)
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
	case !AccessHandler.CheckUserAccess(currentUser, "changepassword"):
		buf = genDataPacket("", "CommandDeniedUserNotAuthorized", status, packet.Saviour.Username)
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
	case !AccessHandler.CheckUserAccess(currentUser, "removeuser"):
		Logger("CommandDeniedUserNotAuthorized", SYSTEM, MSG)
		buf = genDataPacket("", "CommandDeniedUserNotAuhorized", 400, currentUser.GetName())
	default:
		result := CommandHandler.RemoveUser(packet.Saviour.Message, currentUser)
		buf = genDataPacket(currentUser.GetToken(), result, 200, packet.Saviour.Username)
	}
	return buf
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
