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
	permMap   map[string]*Perm
}

type Perm struct {
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	Module string   `json:"module"`
}

// InitAccess is the initializing function that loads the AccessHandler global variable
func InitAccess() {
	var access Access
	access.roleIDMap = access.GetRoleIDMap()
	access.loadDefaultAccess()
	AccessHandler = access
}

// createRole
func (access Access) createRole(roleName string) {
	Logger("CreatingNewRole", ACCESS, MSG)
	createRole := DBHandler.SetupExec("INSERT INTO roles ( name ) VALUES ( ? )", roleName)
	DBHandler.Exec(createRole)
	access.roleIDMap = access.GetRoleIDMap()
}

// cloneRole
func (access Access) cloneRole(roleName string, cloneRole string) {
	access.createRole(roleName)
	for _, perm := range access.permMap {
		Logger("CheckingClonedPerm::"+perm.Name, ACCESS, MSG)
		for _, role := range perm.Roles {
			if role == cloneRole {
				Logger("CloningRole::"+cloneRole, ACCESS, MSG)
				perm.Roles = append(perm.Roles, roleName)
				access.SetRoleAccess(perm.Module, perm.Name, roleName)
			}
		}
	}
}

// removeRole
func (access Access) removeRole(roleName string) {
	Logger("RemovingRole::"+roleName, ACCESS, MSG)
	rid, exists := access.roleIDMap[roleName]
	switch {
	case roleName == "administrator" || roleName == "authorized.user" || roleName == "unauthorized.user":
		Logger("DefaultRoleRemoval::NotAuthorized", ACCESS, MSG)
	case !exists:
		Logger("RoleDoesntExists", ACCESS, MSG)
	default:
		delRoleFromRoles := DBHandler.SetupExec(`DELETE FROM roles WHERE rid = ?`, rid)
		delRoleFromUserRoles := DBHandler.SetupExec(`DELETE FROM user_roles WHERE rid = ?`, rid)
		delRoleFromPermissions := DBHandler.SetupExec(`DELETE FROM role_permissions WHERE rid = ?`, rid)
		DBHandler.Exec(delRoleFromRoles, delRoleFromUserRoles, delRoleFromPermissions)
		for _, perm := range access.permMap {
			for addr, value := range perm.Roles {
				if value == roleName {
					perm.Roles = append(perm.Roles[:addr], perm.Roles[addr+1:]...)
				}
			}
		}
	}
}

// loadDefaultAccess
func (access Access) loadDefaultAccess() {
	var err error
	var tableCount int
	access.permMap = make(map[string]*Perm)
	err = DBHandler.Query(`SELECT COUNT(*) FROM role_permissions`).Scan(&tableCount)
	if err != nil {
		Logger("TableCountFailed::"+err.Error(), "ACCESS", FATAL)
	}
	if tableCount != 0 {
		_, rows := DBHandler.QueryRows(`SELECT rid, module, permission FROM role_permissions`)
		idRoleMap := access.GetIDRoleMap()
		for rows.Next() {
			var rid int
			var module, permission string
			rows.Scan(&rid, &module, &permission)
			_, exists := access.permMap[permission]
			if exists {
				access.permMap[permission].Roles = append(access.permMap[permission].Roles, idRoleMap[rid])
			} else {
				newPerm := Perm{
					Name:   permission,
					Roles:  []string{idRoleMap[rid]},
					Module: module,
				}
				Logger("Loading::Perm::"+permission+"::"+module+"::"+idRoleMap[rid], ACCESS, MSG)
				access.permMap[permission] = &newPerm
			}
		}
	} else {
		permsMap := make(map[string]*Perm)
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
			access.SetRoleAccess(perm.Module, perm.Name, perm.Roles...)
			access.permMap[perm.Name] = perm
		}
	}
}

// AddPermission will add a new permission unless it already exists
func (access Access) AddPermission(perm Perm) {
	_, exists := access.permMap[perm.Name]
	if !exists {
		access.permMap[perm.Name] = &perm
		access.SetRoleAccess(perm.Module, perm.Name, perm.Roles...)
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
			// Do Nothing
		default:
			allowed = true
		}
	}
	return allowed
}

func (access Access) RemoveRoleAccess(role string, module string, permission string) {
	rid := access.roleIDMap[role]
	deleteRolePermission := DBHandler.SetupExec(`DELETE FROM role_permissions `+
		`WHERE rid = ? AND module = ? AND permission = ?`, rid, module, permission)
	DBHandler.Exec(deleteRolePermission)
	Logger("RemovedRolePermission::"+role+"::"+module+"::"+permission, "ACCESS", MSG)
}

func (access Access) SetRoleAccess(module string, permission string, roles ...string) {
	for _, role := range roles {
		rid := access.roleIDMap[role]
		insertRoleAccess := DBHandler.SetupExec(`INSERT INTO role_permissions (rid, module, permission) `+
			`VALUES (?, ?, ?)`, rid, module, permission)
		Logger("WritingRolePermission::"+role+"::"+module+"::"+permission, "ACCESS", MSG)
		DBHandler.Exec(insertRoleAccess)
	}
}

func (access Access) GetIDRoleMap() map[int]string {
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
