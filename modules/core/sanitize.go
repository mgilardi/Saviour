package core

import (
	"errors"
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

const (
	thisModuleSanitize = "Sanitize"
)

// sanitizePacket validates regular data packet if fail sanitize
func sanitizePacket(packet DataPacket) DataPacket {
	validate := validator.New()
	err := validate.Struct(packet)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "User":
				LogHandler.Warn(errors.New("PacketValidation::UserField"), thisModuleSanitize)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.User = regex.ReplaceAllString(packet.Login.User, "")
				packet.Login.User = checkDataSize(packet.Login.User, 45)
			case "Pass":
				LogHandler.Warn(errors.New("PacketValidation::PassField"), thisModuleSanitize)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Login.Pass = regex.ReplaceAllString(packet.Login.Pass, "")
				packet.Login.Pass = checkDataSize(packet.Login.Pass, 45)
			case "Email":
				LogHandler.Warn(errors.New("PacketValidation::EmailField"), thisModuleSanitize)
				regex := regexp.MustCompile("[^a-zA-Z0-9._@]")
				packet.Login.Email = regex.ReplaceAllString(packet.Login.Email, "")
				packet.Login.Email = checkDataSize(packet.Login.Email, 45)
			case "Username":
				LogHandler.Warn(errors.New("PacketValidation::UsernameField"), thisModuleSanitize)
				regex := regexp.MustCompile("[^a-zA-Z0-9]")
				packet.Saviour.Username = regex.ReplaceAllString(packet.Saviour.Username, "")
				packet.Saviour.Username = checkDataSize(packet.Saviour.Username, 45)
			case "Token":
				LogHandler.Warn(errors.New("PacketValidation::TokenField"), thisModuleSanitize)
				regex := regexp.MustCompile(`[^A-Za-z0-9+-\/=]`)
				packet.Saviour.Token = regex.ReplaceAllString(packet.Saviour.Token, "")
				packet.Saviour.Token = checkDataSize(packet.Saviour.Token, 45)
			case "Message":
				LogHandler.Warn(errors.New("PacketValidation::MessageField"), thisModuleSanitize)
				regex := regexp.MustCompile(`[^A-Za-z0-9+-\/=]`)
				packet.Saviour.Message = regex.ReplaceAllString(packet.Saviour.Message, "")
				packet.Saviour.Message = checkDataSize(packet.Saviour.Message, 45)
			case "Status":
				LogHandler.Warn(errors.New("PacketValidation::StatusField"), thisModuleSanitize)
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
