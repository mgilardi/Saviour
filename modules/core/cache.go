package core

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"strconv"
	"time"
)

const (
	// MODULECACHE contains the name for the modules cache
	MODULECACHE = "Cache"
)

// CacheHandler Global Access
var CacheHandler *Cache

// CacheObj interface for any struct containing the Cache method
type CacheObj interface {
	Cache() (string, map[string]interface{})
	CacheID() string
}

// CacheData binary encoding struct for the cache map
type CacheData struct {
	DataMap map[string]interface{}
}

// Cache struct holds the
type Cache struct {
	dbExpireTime time.Duration
	memCache     map[string]map[string]interface{}
}

// InitCache sets up the CacheHandler global variable
func InitCache() {
	var cache Cache
	cache.memCache = make(map[string]map[string]interface{})
	options := OptionsHandler.GetOption("Core")
	cache.dbExpireTime = time.Duration(int(options["DBExpireTime"].(float64))) * time.Minute
	CacheHandler = &cache
	CronHandler.Add(func() {
		CacheHandler.CheckCache()
	})
}

// Cache is the function that is called to load the cache
func (cache *Cache) Cache(obj CacheObj) (bool, map[string]interface{}) {
	var exists bool
	var cacheMap map[string]interface{}
	Logger("CheckingDBCache", PACKAGE+"."+MODULECACHE+".Cache", MSG)
	exists, cacheMap = cache.GetCache(obj)
	if !exists {
		Logger("CouldNotLoadCacheMap", PACKAGE+"."+MODULECACHE+".Cache", ERROR)
	}
	Logger("DBCacheFound", PACKAGE+"."+MODULECACHE+".Cache", MSG)
	return exists, cacheMap
}

// Update will update the cache objects
func (cache *Cache) Update(obj CacheObj) {
	cache.DeleteDBCache(obj)
}

// DeleteDBCache deletes cache entrys for refresh
func (cache *Cache) DeleteDBCache(obj CacheObj) {
	cid := obj.CacheID()
	deleteFromCache := DBHandler.SetupExec(`DELETE FROM cache WHERE cid = ?`, cid)
	DBHandler.Exec(deleteFromCache)
}

// GetCache will return the cache map if the map is not in the cache it will
// be automatically loaded
func (cache *Cache) GetCache(obj CacheObj) (bool, map[string]interface{}) {
	var cacheData CacheData
	var dbData []byte
	var expires sql.NullInt64
	var exists bool
	var buf bytes.Buffer
	var cid string
	var err error
	exists = false
	cid = obj.CacheID()
	err = DBHandler.sql.QueryRow(
		`SELECT data, expires FROM cache `+
			`WHERE cid = ?`, cid).Scan(&dbData, &expires)
	switch {
	case err == sql.ErrNoRows:
		// Cache Row Not found
		Logger("CacheDBEntryNotFound::Creating", PACKAGE+"."+MODULECACHE+".GetCache", MSG)
		cache.SetCache(obj)
		DBHandler.sql.QueryRow(
			`SELECT data, expires FROM cache `+
				`WHERE cid = ?`, cid).Scan(&dbData, &expires)
		exists = true
	case err != nil:
		Logger(err.Error(), PACKAGE+"."+MODULECACHE+".GetCache", ERROR)
	case expires.Valid:
		if !(expires.Int64 < time.Now().Unix()) {
			exists = true
		} else {
			Logger("CacheDBEntryExpired::Creating", PACKAGE+"."+MODULECACHE+".GetCache", MSG)
			cache.SetCache(obj)
			DBHandler.sql.QueryRow(
				`SELECT data, expires FROM cache `+
					`WHERE cid = ?`, cid).Scan(&dbData, &expires)
			exists = true
		}
	default:
		exists = true
	}
	if exists {
		gob.Register(CacheData{})
		decoder := gob.NewDecoder(&buf)

		_, err = buf.Write(dbData)
		if err != nil {
			Logger(err.Error(), PACKAGE+"."+MODULECACHE+".GetCache", ERROR)
		}
		err := decoder.Decode(&cacheData)
		if err != nil {
			Logger(err.Error(), PACKAGE+"."+MODULECACHE+".GetCache", ERROR)
		}
	}
	return exists, cacheData.DataMap
}

// SetCache makes a new cache entry in the database
func (cache *Cache) SetCache(obj CacheObj) {
	var cacheData CacheData
	var buf bytes.Buffer
	var cid string
	var err error
	cid, cacheData.DataMap = obj.Cache()
	gob.Register(CacheData{})
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(&cacheData)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULECACHE+".SetCache", ERROR)
	}
	currentTime := time.Now().Unix()
	expTime := time.Now().Add(cache.dbExpireTime).Unix()
	Logger("WritingCache::"+cid+":"+strconv.FormatInt(expTime, 10), PACKAGE+"."+MODULECACHE+".SetCache", MSG)
	insertCache := DBHandler.SetupExec(
		`INSERT INTO cache (cid, data, created, expires) `+
			`VALUES (?, ?, ?, ?) ON DUPLICATE KEY `+
			`UPDATE data = ?, expires = ?`, cid, buf.Bytes(), currentTime, expTime, buf.Bytes(), expTime)
	DBHandler.Exec(insertCache)
}

// ClearCache clears all records in cache table
func (cache *Cache) ClearCache() {
	truncateTable := DBHandler.SetupExec(`TRUNCATE TABLE cache`)
	DBHandler.Exec(truncateTable)
}

// CheckCache iterates through cache records checks if the record is expired
// if the record is expired the record is deleted, if the record has NULL for
// a value it is permanent and is skipped, if the record has not expired it is skipped
func (cache *Cache) CheckCache() {
	var cid string
	var expires sql.NullInt64
	rows, err := DBHandler.sql.Query(`SELECT cid, expires FROM cache`)
	Logger("CheckDBCacheForExpired", PACKAGE+"."+MODULECACHE+".CheckCache", MSG)
	switch {
	case err == sql.ErrNoRows:
		Logger("CacheTableEmpty", PACKAGE+"."+MODULECACHE+".CheckCache", MSG)
	case err != nil && err.Error() != "EOF":
		Logger(err.Error(), PACKAGE+"."+MODULECACHE+".CheckCache", ERROR)
	default:
		for rows.Next() {
			rows.Scan(&cid, &expires)
			Logger("DBCacheEntry::"+cid+"::Expires::"+strconv.FormatInt(expires.Int64, 10)+
				"::TimeNow::"+strconv.FormatInt(time.Now().Unix(), 10), PACKAGE+"."+MODULECACHE+".CheckCache", MSG)
			if expires.Valid && expires.Int64 < time.Now().Unix() {
				Logger("RemovingExpired::"+cid, PACKAGE+"."+MODULECACHE+".CheckCache", MSG)
				if err != nil {
					Logger(err.Error(), PACKAGE+"."+MODULECACHE+".CheckCache", ERROR)
				} else {
					deleteFromCache := DBHandler.SetupExec(`DELETE FROM cache WHERE cid = ?`, cid)
					DBHandler.Exec(deleteFromCache)
				}
			}
		}
	}
}
