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
	dbExpireTime  time.Duration
	memExpireTime time.Duration
	memCache      map[string]map[string]interface{}
}

// InitCache sets up the CacheHandler global variable
func InitCache() {
	var cache Cache
	cache.memCache = make(map[string]map[string]interface{})
	options := OptionsHandler.GetOptions("core")
	cache.dbExpireTime = time.Duration(int(options["DBExpireTime"].(float64))) * time.Minute
	cache.memExpireTime = time.Duration(int(options["MemExpireTime"].(float64))) * time.Minute
	CacheHandler = &cache
	OptionsHandler.CacheLoaded()
	CronHandler.Add(func() {
		CacheHandler.CheckCache()
		CacheHandler.CheckMemCache()
	})
}

// Cache is the function that is called to load the cache
func (cache *Cache) Cache(obj CacheObj) (bool, map[string]interface{}) {
	var exists bool
	var cacheMap map[string]interface{}
	Logger("CheckingMemCache", PACKAGE+"."+MODULECACHE+".Cache", MSG)
	exists, cacheMap = cache.GetMemCache(obj)
	if exists {
		Logger("ReceivedMemCache", PACKAGE+"."+MODULECACHE+".Cache", MSG)
		return exists, cacheMap
	}
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
	cache.DeleteMemCache(obj)
	cache.DeleteDBCache(obj)
}

// GetMemCache checks for memory cache and returns value if exists if not
// it will load a new entry into the memory cache
func (cache *Cache) GetMemCache(obj CacheObj) (bool, map[string]interface{}) {
	value, exists := cache.memCache[obj.CacheID()]
	if !exists {
		Logger("CacheMemEntryNotFound::Creating", PACKAGE+"."+MODULECACHE+".GetMemCache", MSG)
		cache.SetMemCache(obj)
	} else {
		Logger("MemoryCacheEntryExists", PACKAGE+"."+MODULECACHE+".GetMemCache", MSG)
		if time.Now().Unix() > value["expires"].(int64) {
			exists = false
			Logger("CacheMemEntryExpired::Creating", PACKAGE+"."+MODULECACHE+".GetMemCache", MSG)
			cache.SetMemCache(obj)
		}
	}
	return exists, value
}

// SetMemCache adds a new entry into the memory cache
func (cache *Cache) SetMemCache(obj CacheObj) {
	cid, cacheMap := obj.Cache()
	cacheMap["expires"] = time.Now().Add(cache.memExpireTime).Unix()
	cache.memCache[cid] = cacheMap

}

// CheckMemCache removes expired entrys from the memory cache
func (cache *Cache) CheckMemCache() {
	Logger("CheckMemCacheForExpired", PACKAGE+"."+MODULECACHE+".CheckMemCache", MSG)
	for key, value := range cache.memCache {
		Logger("MemCacheEntry::"+key+"::Expires::"+strconv.FormatInt(value["expires"].(int64), 10)+
			"::TimeNow::"+strconv.FormatInt(time.Now().Unix(), 10), PACKAGE+"."+MODULECACHE+".CheckMemCache", MSG)
		if time.Now().Unix() > value["expires"].(int64) {
			Logger("RemovingExpired::"+key, PACKAGE+"."+MODULECACHE+".CheckMemCache", MSG)
			delete(cache.memCache, key)
		}
	}
}

// DeleteMemCache deletes cache entrys for refresh
func (cache *Cache) DeleteMemCache(obj CacheObj) {
	cid := obj.CacheID()
	for k := range cache.memCache {
		if cid == k {
			delete(cache.memCache, k)
		}
	}
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
	var expTime int64
	cid, cacheData.DataMap = obj.Cache()
	gob.Register(CacheData{})
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(&cacheData)
	if err != nil {
		Logger(err.Error(), PACKAGE+"."+MODULECACHE+".SetCache", ERROR)
	}
	expTime = time.Now().Add(cache.dbExpireTime).Unix()
	Logger("WritingCache::"+cid+":"+strconv.FormatInt(expTime, 10), PACKAGE+"."+MODULECACHE+".SetCache", MSG)
	insertCache := DBHandler.SetupExec(
		`INSERT INTO cache (cid, data, expires) `+
			`VALUES (?, ?, ?) ON DUPLICATE KEY `+
			`UPDATE data = ?, expires = ?`, cid, buf.Bytes(), expTime, buf.Bytes(), expTime)
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
