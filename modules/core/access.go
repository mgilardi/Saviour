package core

import (
	"database/sql"
)

const (
	// ACCESS constant for module output of logger
	ACCESS = "Access"
)

// AccessHandler is the global variable to access the access module
var AccessHandler *Access

// AccessObj is a container for the permissions for each loaded module
type AccessObj struct {
	name         string
	allowedRoles map[string]bool
}

func makeAccessObj(name string) *AccessObj {
	var obj AccessObj
	obj.allowedRoles = make(map[string]bool)
	obj.name = name
	return &obj
}

// Perm interface for any module that requires access control
type Perm interface {
	PermID() string
	DefaultPerm() map[string]bool
}

// Access struct contains the map of permissions loaded for each module
type Access struct {
	accessMap map[string]*AccessObj
}

// InitAccess is the initial configuration point for the access module
func InitAccess() {
	var access Access
	Logger("StartingAccess", ACCESS, MSG)
	access.accessMap = make(map[string]*AccessObj)
	access.loadDB()
	AccessHandler = &access
}

// CreateUserRole creates a new user role named with input
func (access *Access) CreateUserRole(roleName string) {
	Logger("CreatingUserRole", ACCESS, MSG)
	_, err := DBHandler.sql.Exec(`INSERT INTO role (name) VALUES (?)`, roleName)
	if err != nil {
		Logger(err.Error(), ACCESS, ERROR)
	}
	access.reloadPerm()
	access.updateDB()
}

// RemoveUserRole removes a user role that matches the input name
func (access *Access) RemoveUserRole(roleName string) {
	Logger("RemovingUserRole", ACCESS, MSG)
	roleMap := access.genRoleNameMap()
	tx, err := DBHandler.sql.Begin()
	_, err = tx.Exec(`DELETE FROM role WHERE name = ?`, roleName)
	_, err = tx.Exec(`DELETE FROM user_roles WHERE rid = ?`, roleMap[roleName])
	if err != nil {
		tx.Rollback()
		Logger(err.Error(), ACCESS, ERROR)
	} else {
		tx.Commit()
		access.clearDB()
		access.reloadPerm()
		access.updateDB()
	}
}

// GetUserRoles returns a map containing the user roles and there role id
func (access *Access) GetUserRoles() map[string]int {
	return access.genRoleNameMap()
}

// AllowAccess will disable access for a usertype
func (access *Access) AllowAccess(usrTyp string, obj Perm) {
	access.accessMap[obj.PermID()].allowedRoles[usrTyp] = true
	access.updateDB()
}

// DisableAccess will disable access for a usertype
func (access *Access) DisableAccess(usrTyp string, obj Perm) {
	access.accessMap[obj.PermID()].allowedRoles[usrTyp] = false
	access.updateDB()
}

// GetPermMap will return the permissions map for input module
func (access *Access) GetPermMap(obj Perm) map[string]bool {
	return access.accessMap[obj.PermID()].allowedRoles
}

// GetAccessMap returns the accessMap which contains all loaded modules
func (access *Access) GetAccessMap() map[string]*AccessObj {
	return access.accessMap
}

// CheckPerm check if a user has access to a module
func (access *Access) CheckPerm(user *User, obj Perm) bool {
	valid := false
	for name := range user.GetRoleNames() {
		if access.accessMap[obj.PermID()].allowedRoles[name] {
			valid = true
			break
		}
	}
	return valid
}

// LoadPerm loads the permissions from the provided module that
// matches the Perm interface
func (access *Access) LoadPerm(obj Perm) {
	_, exists := access.accessMap[obj.PermID()]
	if !exists {
		access.accessMap[obj.PermID()] = makeAccessObj(obj.PermID())
		usrDefault := obj.DefaultPerm()
		roleMap := access.getRoles()
		for k, v := range usrDefault {
			roleMap[k] = v
		}
		access.accessMap[obj.PermID()].allowedRoles = roleMap
		access.updateDB()
	}
}

func (access *Access) reloadPerm() {
	roleMap := access.getRoles()
	for module, permObj := range access.accessMap {
		for userRole, perm := range permObj.allowedRoles {
			roleMap[userRole] = perm
		}
		access.accessMap[module].allowedRoles = roleMap
	}
}

func (access *Access) getRoles() map[string]bool {
	Logger("GettingRolesMap", ACCESS, MSG)
	rolesMap := make(map[string]bool)
	rows, err := DBHandler.sql.Query(`SELECT name FROM role`)
	if err != nil {
		Logger(err.Error(), ACCESS, ERROR)
	}
	for rows.Next() {
		var role string
		err = rows.Scan(&role)
		if err != nil {
			Logger(err.Error(), ACCESS, ERROR)
		}
		Logger("Loading::"+role, ACCESS, MSG)
		rolesMap[role] = false
	}
	return rolesMap
}

func (access *Access) loadDB() {
	var rid int
	var module string
	var allowed bool
	var name string
	Logger("LoadingAccessFromDB", ACCESS, MSG)
	rows, _ := DBHandler.sql.Query(
		`SELECT user_permissions.rid, module, allowed, name FROM user_permissions ` +
			`INNER JOIN role ON user_permissions.rid = role.rid`)
	for rows.Next() {
		err := rows.Scan(&rid, &module, &allowed, &name)
		if err != nil {
			Logger(err.Error()+"::loadDB::Scan", ACCESS, ERROR)
		}
		_, exists := access.accessMap[module]
		if !exists {
			access.accessMap[module] = makeAccessObj(module)
			access.accessMap[module].allowedRoles[name] = allowed
		} else {
			access.accessMap[module].allowedRoles[name] = allowed
		}
	}
	rows.Close()
	Logger("LoadDBComplete", ACCESS, MSG)
}

func (access *Access) genRoleNameMap() map[string]int {
	var rid int
	var name string
	Logger("GeneratingRoleMap", ACCESS, MSG)
	roleMap := make(map[string]int)
	rows, err := DBHandler.sql.Query(`SELECT rid, name FROM role`)
	if err != nil {
		Logger(err.Error(), ACCESS, ERROR)
	}
	for rows.Next() {
		rows.Scan(&rid, &name)
		roleMap[name] = rid
	}
	return roleMap
}

func (access *Access) updateDB() {
	var err error
	roleMap := access.genRoleNameMap()
	for name, acMap := range access.accessMap {
		Logger("LoadingAccessDB::"+name, ACCESS, MSG)
		for usrTyp, perm := range acMap.allowedRoles {
			var allowed sql.NullBool
			err = DBHandler.sql.QueryRow(
				`SELECT allowed FROM user_permissions `+
					`WHERE rid = ? AND module = ?`, roleMap[usrTyp], name).Scan(&allowed)
			switch {
			case err == sql.ErrNoRows:
				Logger("CreatingNewEntry", ACCESS, MSG)
				_, err = DBHandler.sql.Exec(
					`INSERT INTO user_permissions (rid, module, allowed) `+
						`VALUES (?, ?, ?)`, roleMap[usrTyp], name, perm)
				if err != nil {
					Logger(err.Error(), ACCESS, ERROR)
				}
			case err != nil:
				Logger(err.Error(), ACCESS, ERROR)
			case allowed.Valid && allowed.Bool == perm:
				// Ignore
			default:
				_, err = DBHandler.sql.Exec(
					`UPDATE user_permissions SET allowed = ? `+
						`WHERE rid = ? AND module =?`, perm, roleMap[usrTyp], name)
				if err != nil {
					Logger(err.Error(), ACCESS, ERROR)
				}
			}
		}
	}
}

func (access *Access) clearDB() {
	Logger("ClearingAccessDB", ACCESS, MSG)
	_, err := DBHandler.sql.Exec(`TRUNCATE TABLE user_permissions`)
	if err != nil {
		Logger(err.Error(), ACCESS, MSG)
	}
}
