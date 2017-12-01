package core

import (
	"encoding/json"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
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
		Status   int    `json:"status,omitempty" validate:"lte=600|len=0"`
		Token    string `json:"token,omitempty" validate:"max=45,ascii|len=0"`
		Message  string `json:"message,omitempty" validate:"max=45|len=0"`
	} `json:"saviour"`
}

// genDataPacket generates a packet for transaction
func genDataPacket(token string, message string, status int, username string) []byte {
	var packet DataPacket
	var buf []byte
	// Testing
	packet.Saviour.Token = token
	packet.Saviour.Message = message
	packet.Saviour.Status = status
	packet.Saviour.Username = username
	buf, err := json.Marshal(&packet)
	if err != nil {
		Logger(err.Error(), "Packet", ERROR)
	}
	return buf
}

// loadDataPacket loads incoming packet for analysis
func loadDataPacket(buf []byte) (bool, DataPacket) {
	var packet DataPacket
	valid := true
	err := json.Unmarshal(buf, &packet)
	if err != nil {
		Logger(err.Error(), "Packet", ERROR)
		valid = false
	}
	return valid, sanitizePacket(packet)
}

// sanitizePacket validates regular data packet if fail sanitize
func sanitizePacket(packet DataPacket) DataPacket {
	validate := validator.New()
	err := validate.Struct(packet)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "User":
				Logger("PacketValidation::UserField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Login.User, "Sanitize", MSG)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.User = regex.ReplaceAllString(packet.Login.User, "")
				packet.Login.User = checkDataSize(packet.Login.User, 45)
				Logger("Sanitized::"+packet.Login.User, "Sanitize", MSG)
			case "Pass":
				Logger("PacketValidation::PassField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Login.Pass, "Sanitize", MSG)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.Pass = regex.ReplaceAllString(packet.Login.Pass, "")
				packet.Login.Pass = checkDataSize(packet.Login.Pass, 45)
				Logger("Sanitized::"+packet.Login.Pass, "Sanitize", MSG)
			case "Email":
				Logger("PacketValidation::EmailField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Login.Email, "Sanitize", MSG)
				regex := regexp.MustCompile("[^a-zA-Z0-9._@]")
				packet.Login.Email = regex.ReplaceAllString(packet.Login.Email, "")
				packet.Login.Email = checkDataSize(packet.Login.Email, 45)
				Logger("Sanitized::"+packet.Login.Email, "Sanitize", MSG)
			case "Username":
				Logger("PacketValidation::UsernameField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Saviour.Username, "Sanitize", MSG)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Saviour.Username = regex.ReplaceAllString(packet.Saviour.Username, "")
				packet.Saviour.Username = checkDataSize(packet.Saviour.Username, 45)
				Logger("Sanitized::"+packet.Saviour.Username, "Sanitize", MSG)
			case "Token":
				Logger("PacketValidation::TokenField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Saviour.Token, "Sanitize", MSG)
				regex := regexp.MustCompile(`[^A-Za-z0-9+-_\/=]`)
				packet.Saviour.Token = regex.ReplaceAllString(packet.Saviour.Token, "")
				packet.Saviour.Token = checkDataSize(packet.Saviour.Token, 45)
				Logger("Sanitized::"+packet.Saviour.Token, "Sanitize", MSG)
			case "Message":
				Logger("PacketValidation::MessageField", "Sanitize", WARN)
				Logger("Sanitizing::"+packet.Saviour.Message, "Sanitize", MSG)
				regex := regexp.MustCompile(`[^A-Za-z0-9+-:_\/=]`)
				packet.Saviour.Message = regex.ReplaceAllString(packet.Saviour.Message, "")
				packet.Saviour.Message = checkDataSize(packet.Saviour.Message, 45)
				Logger("Sanitized::"+packet.Saviour.Message, "Sanitize", MSG)
			case "Status":
				Logger("PacketValidation::StatusField", "Sanitize", WARN)
				packet.Saviour.Status = 200
			default:
				// Ignore
			}
		}
	}
	return packet
}

// checkDataSize trims down strings if they are above validation size
func checkDataSize(data string, size int) string {
	var trimStr string
	if len(data) > size {
		trimStr = data[0:size]
	} else {
		trimStr = data
	}
	return trimStr
}
