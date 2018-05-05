package core

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
)

const (
	// ACCESS constant for module output of logger
	ACCESS = "Access"
)

var AccessHandler Access

type Access struct {
	roleIDMap map[string]int
}

type Permission struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Module string `json:"module"`
}

func InitAccess() {
	var access Access
	access.roleIDMap = access.GetRoleIDMap()
	access.loadDefaultAccess()
	AccessHandler = access
}

func (access Access) loadDefaultAccess() {
	var err error
	var tableCount int
	err = DBHandler.Query(`SELECT COUNT(*) FROM role_permissions`).Scan(&tableCount)
	if err != nil {
		Logger("TableCountFailed::"+err.Error(), "ACCESS", FATAL)
	}
	if tableCount != 0 {
		// Do Nothing
	} else {
		permsMap := make(map[string]Permission)
		path, err := FindPath("config")
		if err != nil {
			Logger("DefaultAccessPermsFailed::"+err.Error(), "ACCESS", FATAL)
		}
		raw, err := ioutil.ReadFile(path + "access.json")
		if err != nil {
			Logger("DefaultAccessPermsFailed::"+err.Error(), "ACCESS", FATAL)
		}
		err = json.Unmarshal(raw, &permsMap)
		if err != nil {
			Logger("DefaultAccessPermsFailed::"+err.Error(), "ACCESS", FATAL)
		}
		for _, perm := range permsMap {
			access.SetRoleAccess(perm.Role, perm.Module, perm.Name)
		}
	}
}

func (access Access) GetUserAccess(user *User, permission string) bool {
	allowed := false
	userRoleMap := user.GetUserRoleMap()
	for _, value := range userRoleMap {
		err := DBHandler.Query(`SELECT rid, permission FROM role_permissions `+
			`WHERE rid = ? AND permission = ?`, value, permission).Scan()
		switch {
		case err == sql.ErrNoRows:
			//
		default:
			allowed = true
		}
	}
	return allowed
}

func (access Access) SetRoleAccess(role string, module string, permission string) {
	rid := access.roleIDMap[role]
	insertRoleAccess := DBHandler.SetupExec(`INSERT INTO role_permissions (rid, module, permission) `+
		`VALUES (?, ?, ?)`, rid, module, permission)
	Logger("WritingRolePermission::"+role+"::"+module+"::"+permission, "ACCESS", MSG)
	DBHandler.Exec(insertRoleAccess)
}

func (access Access) GetRoleIDMap() map[string]int {
	roleIDMap := make(map[string]int)
	exists, rows := DBHandler.QueryRows(`SELECT rid, name FROM roles`)
	if !exists {
		Logger("RoleIDMapNotFound", "ACCESS", FATAL)
	}
	defer rows.Close()
	for rows.Next() {
		var rid int
		var name string
		err := rows.Scan(&rid, &name)
		if err != nil {
			Logger("RoleIDMapFailed::"+err.Error(), "ACCESS", FATAL)
		}
		Logger("RoleIDFound::"+name, "ACCESS", MSG)
		roleIDMap[name] = rid
	}
	return roleIDMap
}
