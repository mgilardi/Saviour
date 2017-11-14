PACKAGE DOCUMENTATION

package system
    import "modules/system"


FUNCTIONS

func InitSystem(datab *database.Database, cache *database.Cache)
    InitSystem initialize system

TYPES

type DataPacket struct {
    Login struct {
        User  string `json:"user,omitempty" validate:"min=0,max=45,alphanum"`
        Pass  string `json:"pass,omitempty" validate:"min=0,max=45,alphanum"`
        Email string `json:"email,omitempty" validate:"min=0,max=45,email"`
    } `json:"login"`
    Saviour struct {
        Username string `json:"username,omitempty" validate:"max=45,alphanum"`
        Status   int    `json:"status,omitempty" validate:"max=3"`
        Token    string `json:"token,omitempty" validate:"max=45,base64"`
        Message  string `json:"message,omitempty" validate:"max=45,base64"`
    } `json:"saviour"`
}
    DataPacket is the struct that json files are loaded into when marshaled

type System struct {
    // contains filtered or unexported fields
}
    System contains server responses to http requests
