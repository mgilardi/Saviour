package system

import (
	//"modules/logger"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

// sanitizeLogin validates login packet if fail sanitize
func sanitizeLogin(packet DataPacket) [3]string {
	validate := validator.New()
	err := validate.Struct(packet)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "User":
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.User = regex.ReplaceAllString(packet.Login.User, "")
				packet.Login.User = checkDataSize(packet.Login.User, 45)
			case "Pass":
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.Pass = regex.ReplaceAllString(packet.Login.Pass, "")
				packet.Login.Pass = checkDataSize(packet.Login.Pass, 45)
			case "Email":
				regex := regexp.MustCompile("[^a-zA-Z0-9._@]")
				packet.Login.Email = regex.ReplaceAllString(packet.Login.Email, "")
				packet.Login.Email = checkDataSize(packet.Login.Email, 45)
			default:
				// Ignore
			}

		}
	}
	login := [3]string{packet.Login.User, packet.Login.Pass, packet.Login.Email}

	return login
}

// sanitizePacket validates regular data packet if fail sanitize
func sanitizePacket(packet DataPacket) DataPacket {
	validate := validator.New()
	err := validate.Struct(packet)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Saviour.Username = regex.ReplaceAllString(packet.Saviour.Username, "")
				packet.Saviour.Username = checkDataSize(packet.Saviour.Username, 45)
			case "Token":
				regex := regexp.MustCompile(`[^A-Za-z0-9+-\/=]`)
				packet.Saviour.Token = regex.ReplaceAllString(packet.Saviour.Token, "")
				packet.Saviour.Token = checkDataSize(packet.Saviour.Token, 45)
			case "Message":
				regex := regexp.MustCompile(`[^A-Za-z0-9+-\/=]`)
				packet.Saviour.Message = regex.ReplaceAllString(packet.Saviour.Message, "")
				packet.Saviour.Message = checkDataSize(packet.Saviour.Message, 45)
			case "Status":
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
