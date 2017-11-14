PACKAGE DOCUMENTATION

package user
    import "modules/user"


TYPES

type User struct {
    // contains filtered or unexported fields
}
    User handles users

func InitUser(id int, db *database.Database, cache *database.Cache) *User
    InitUser constructs the user on initial login

func (user *User) CheckToken()
    CheckToken checks to see if a token exists in the database if not it
    generates one.

func (user *User) GetInfoMap() map[string]interface{}
    GetInfoMap retrieves the loaded infomap from the database

func (user *User) GetName() string
    GetName returns username

func (user *User) GetToken() string
    GetToken returns user token

func (user *User) InfoUpdate()
    InfoUpdate calls for the username and token from the database to be held
    in memory

func (user *User) IsOnline() bool
    IsOnline returns the online flag for the user

func (user *User) SetOnline(isOnline bool)
    SetOnline will set the flag to the input

func (user *User) SetToken()
    SetToken generates a new token and writes it to DB
