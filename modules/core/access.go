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
	permMap map[string]*Perm
	roleMap map[string]*int
}

type Perm struct {
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	Module string   `json:"module"`
}

// InitAccess this is the inital loading of the Access Handler, it checks if the
// table has permissions entrys in it, if it does it will trigger a loading of
// the table from the database, if it doesn't have any entrys in the Permissions
// table it will load the default permissions JSON file.
func InitAccess() {
	var err error
	var tableCount int
	var access Access
	access.permMap = make(map[string]*Perm)
	access.roleMap = make(map[string]*int)
	err = DBHandler.Query(`SELECT COUNT(*) FROM role_permissions`).Scan(&tableCount)
	switch {
	case err != nil:
		Logger("TableCountFailed::"+err.Error(), ACCESS, ERROR)
	case tableCount != 0:
		access.ReloadPerms()
	default:
		path, err := FindPath("config")
		if err != nil {
			Logger("FindMsgFailed", ACCESS, ERROR)
		}
		raw, err := ioutil.ReadFile(path + "access.json")
		if err != nil {
			Logger("ReadFileFailed", ACCESS, ERROR)
		}
		err = json.Unmarshal(raw, &access.permMap)
		if err != nil {
			Logger("UnmarshalFailed", ACCESS, ERROR)
		}
		for _, perm := range access.permMap {
			access.UpdatePerms(perm)
		}
	}
	AccessHandler = access
}

// CreateRole Will create a new role using the default user permissions
func (access Access) CreateRole(roleName string) {
	err := DBHandler.Query(`SELECT rid FROM roles WHERE name = ?`, roleName).Scan()
	if err != sql.ErrNoRows {
		Logger("RoleExists", ACCESS, ERROR)
	} else {
		addRole := DBHandler.SetupExec(`INSERT INTO roles (name) VALUES (?)`, roleName)
		DBHandler.Exec(addRole)
		for _, perm := range access.permMap {
			for _, role := range perm.Roles {
				if role == "authorized.user" {
					perm.Roles = append(perm.Roles, roleName)
					access.UpdatePerms(perm)
				}
			}
		}
		Logger("RoleCreated::"+roleName, ACCESS, MSG)
	}
}

// CloneRole will clone the role provided to it
func (access Access) CloneRole(roleName string, cloneName string) {
	var role, clone int
	errRole := DBHandler.Query(`SELECT rid FROM roles WHERE name = ?`, roleName).Scan(&role)
	errClone := DBHandler.Query(`SELECT rid FROM roles WHERE name = ?`, cloneName).Scan(&clone)
	switch {
	case errRole != sql.ErrNoRows:
		Logger("RoleExists", ACCESS, ERROR)
	case errClone == sql.ErrNoRows:
		Logger("CloneRoleCannotBeFound", ACCESS, ERROR)
	default:
		addRole := DBHandler.SetupExec(`INSERT INTO roles (name) VALUES (?)`, roleName)
		DBHandler.Exec(addRole)
		for _, perm := range access.permMap {
			for _, role := range perm.Roles {
				if role == cloneName {
					perm.Roles = append(perm.Roles, roleName)
					access.UpdatePerms(perm)
				}
			}
		}
		Logger("RoleCreated::"+roleName+"::Cloned::"+cloneName, ACCESS, MSG)
	}
}

// RemoveRole
func (access Access) RemoveRole(roleName string) {
	err := DBHandler.Query(`SELECT FROM roles WHERE name = ?`, roleName).Scan()
	switch {
	case roleName == "administrator" || roleName == "authorized.user" || roleName == "unauthorized.user":
		Logger("CannotRemoveDefaultRoles", ACCESS, ERROR)
	case err == sql.ErrNoRows:
		Logger("RemoveRoleDoesNotExist", ACCESS, ERROR)
	default:
		deletePerm := DBHandler.SetupExec(`DELETE FROM role_permissions WHERE rid = ?`, access.roleMap[roleName])
		deleteRole := DBHandler.SetupExec(`DELETE FROM roles WHERE name = ?`, roleName)
		DBHandler.Exec(deleteRole, deletePerm)
		access.ReloadPerms()
	}
}

func (access Access) CheckRole(roleName string) bool {
	var exists bool
	err := DBHandler.Query(`SELECT rid FROM roles WHERE name = ?`, roleName).Scan()
	if err == sql.ErrNoRows {
		Logger("CheckRoleDoesNotExist", ACCESS, ERROR)
		exists = false
	} else {
		exists = true
	}
	return exists
}

func (access Access) GetPerms() map[string]*Perm {
	access.ReloadPerms()
	return access.permMap
}

func (access Access) UpdatePerms(newPerm *Perm) {
	var dbEntry [][]interface{}
	updated := false
	access.genRoleMap()
	for _, role := range newPerm.Roles {
		err := DBHandler.Query(`SELECT FROM role_permissions WHERE rid = ? AND module = ? `+
			`AND permission = ? `, access.roleMap[role], newPerm.Module, newPerm.Name).Scan()
		if err == sql.ErrNoRows {
			updated = true
			dbEntry = append(dbEntry, DBHandler.SetupExec(`INSERT INTO role_permissions (rid, module, permission) `+
				`VALUES (?, ?, ?)`, access.roleMap[role], newPerm.Module, newPerm.Name))
		}
	}
	if updated {
		DBHandler.Exec(dbEntry...)
	}
	access.ReloadPerms()
}

func (access Access) ReloadPerms() {
	_, rows := DBHandler.QueryRows(`SELECT rid, module, permission FROM role_permissions`)
	idRoleMap := access.genIDMap()
	for rows.Next() {
		var rid int
		var module, permission string
		rows.Scan(&rid, &module, &permission)
		_, exists := access.permMap[permission]
		if exists {
			access.permMap[permission].Roles = append(access.permMap[permission].Roles, idRoleMap[rid])
		} else {
			access.permMap[permission] = &Perm{
				Name:   permission,
				Roles:  []string{idRoleMap[rid]},
				Module: module,
			}
		}
	}
}

func (access Access) CheckUserAccess(user *User, permission string) bool {
	allowed := false
	userRoleMap := user.GetUserRoleMap()
	for _, value := range userRoleMap {
		err := DBHandler.Query(`SELECT rid, permission FROM role_permissions `+
			`WHERE rid = ? AND permission = ?`, value, permission).Scan()
		switch {
		case err == sql.ErrNoRows:
			// Do Nothing
		default:
			allowed = true
		}
	}
	return allowed
}

func (access Access) genRoleMap() {
	exists, rows := DBHandler.QueryRows(`SELECT rid, name FROM roles`)
	if !exists {
		// Log
	}
	defer rows.Close()
	for rows.Next() {
		var rid int
		var name string
		err := rows.Scan(&rid, &name)
		if err != nil {
			Logger("RoleIDMapFailed", ACCESS, MSG)
		}
		access.roleMap[name] = &rid
	}
}

// GetIDRoleMap returns a map of the rid idenitifying the role name
func (access Access) genIDMap() map[int]string {
	idRoleMap := make(map[int]string)
	exists, rows := DBHandler.QueryRows(`SELECT rid, name FROM roles`)
	if !exists {
		Logger("IDRoleMapNotFound", ACCESS, FATAL)
	}
	defer rows.Close()
	for rows.Next() {
		var rid int
		var name string
		err := rows.Scan(&rid, &name)
		if err != nil {
			Logger("IDRoleMapFailed"+err.Error(), ACCESS, FATAL)
		}
		idRoleMap[rid] = name
	}
	return idRoleMap
}
