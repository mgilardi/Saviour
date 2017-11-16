PACKAGE DOCUMENTATION

package database
    import "modules/database"

    Package database checks database exists initializes the Database object
    provides object methods for reading/writing/managing of the database.
    the database is limited to pre-defined actions in the form of routines
    for security purposes

TYPES

type Cache struct {
    // contains filtered or unexported fields
}
    Cache holds the options, buffer, and db access

func InitCache(db *Database) *Cache
    InitCache constructs the cache object

func (cache *Cache) CheckCache()
    CheckCache checks for expired cache entrys

func (cache *Cache) ClearAllCache()
    ClearAllCache removes all the records from the cache and reloads module
    options

func (cache *Cache) GetCacheMap(cid string) (bool, map[string]interface{})
    GetCacheMap returns requested cache map

func (cache *Cache) SetCacheMap(cid string, data map[string]interface{}, expires bool)
    SetCacheMap converts map into a binary string and loads the string into
    the database if the expires flag is false then the cached map will be
    added to the table with NULL for the expires tag and will not be removed
    when a CheckCache() occurs. If the expires flag is set to true a unix
    time will be written to the expires column of the record base on the
    configuration value of ExpireTime in minutes.

type Data struct {
    DataMap map[string]interface{}
}
    Data holds the data input/output of the binary encoder/decoder

type Database struct {
    // contains filtered or unexported fields
}
    Database tyoe contains the sql access, options, logger, and the dsn for
    sql login

func InitDatabase() *Database
    InitDatabase initialize the database object and passes a pointer to the
    main loop

func (db *Database) CheckCache()
    CheckCache iterates through cache records checks if the record is
    expired if the record is expired the record is deleted, if the record
    has NULL for a value it is permanent and is skipped, if the record has
    not expired it is skipped

func (db *Database) CheckDB()
    CheckDB checks if database exists and outputs tables that are found.

func (db *Database) CheckToken(uid int) bool
    CheckToken checks if the user has a token in the database login_token
    table

func (db *Database) CheckUserExits(name string) bool
    CheckUserExits checks the database for a username and returns true or
    false if it exists

func (db *Database) CheckUserLogin(name string, pass string) (bool, int)
    CheckUserLogin inputs username and password and checks if it matches the
    database returns a boolean if the user if found along with the users
    database uid

func (db *Database) ClearCache()
    ClearCache clears all records in cache table

func (db *Database) CreateUser(name string, pass string, email string)
    CreateUser creates a new user entry in the database

func (db *Database) GetUserID(name string) (int, error)
    GetUserID will return the database uid for a username

func (db *Database) GetUserMap(uid int) (map[string]interface{}, error)
    GetUserMap loads user information from the database and returns it
    inside of a map

func (db *Database) ReadCache(cid string) []byte
    ReadCache returns a cache entry

func (db *Database) RemoveUser(name string)
    RemoveUser removes a user entry from the database

func (db *Database) StoreToken(uid int, token string)
    StoreToken writes user token to the database

func (db *Database) WriteCache(cid string, data []byte)
    WriteCache creates a new cache entry

func (db *Database) WriteCacheExp(cid string, data []byte, expires int64)
    WriteCacheExp creates a new cache entry that expires
