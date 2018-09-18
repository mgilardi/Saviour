package core

import (
	"testing"
)

func TestPacket_genDataPacket(t *testing.T) {
	exists, userMap, _ := InitUser("Admin", "Password")
	match := `{"login":{},"saviour":{"username":"Admin","status":200,"message":"test"}}`
	if !exists {
		t.Error("AdminUserNotFound")
	} else {
		buf := genDataPacket(userMap["token"].(string), "test", 200, userMap["name"].(string))
		if match != string(buf) {
			t.Error("genDataPacketMisMatch")
		}
	}
}

func TestPacket_loadDataPacket(t *testing.T) {
	var packet DataPacket
	buf := genDataPacket("@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls", "@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls", 6666, "@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls")
	_, packet = loadDataPacket(buf)
	packet.Login.User = "@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls"
	packet.Login.Pass = "@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls"
	packet.Login.Email = "@%^&*flkfalhfdfhklhfklafhklahkldafhiuahjhfldlahdjlafjyuuogihfflfgsljkhgkjshgls"
	sanitizePacket(packet)
}
