package core

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"strings"
	"time"
)

// CacheHandler global cache access
var CacheHandler *Cache

// Data holds the data input/output of the binary encoder/decoder
type Data struct {
	DataMap map[string]interface{}
}

// Cache holds the options, buffer, and db access
type Cache struct {
	expireTime time.Duration
	buf        bytes.Buffer
	db         *Database
}

// InitCache constructs the cache object
func InitCache(db *Database) {
	var newCache Cache
	DebugHandler.Sys("Starting", "Database")
	newCache.db = db
	options := GetOptions("Core")
	newCache.expireTime = time.Duration(int(options["ExpireTime"].(float64)))
	newCache.CacheOptions()
	CacheHandler = &newCache
	CronHandler.Add("CacheCheck", true, func() {
		CacheHandler.CheckCache()
	})
}

// CacheOptions loads the modules configuration files into cache
func (cache *Cache) CacheOptions() {
	DebugHandler.Sys("CacheOptions", "Cache")
	allOptions := GetAllOptions()
	for _, opt := range allOptions {
		cache.SetCacheMap(strings.ToLower(opt["Name"].(string)+":config"), opt, false)
	}
}

// CheckCache checks for expired cache entrys
func (cache *Cache) CheckCache() {
	DebugHandler.Sys("CheckCacheForExpired", "Cache")
	cache.checkCache()
}

// ClearAllCache removes all the records from the cache and reloads module options
func (cache *Cache) ClearAllCache() {
	DebugHandler.Sys("ClearAllCache", "Cache")
	cache.clearCache()
	//cache.CacheOptions()
}

// SetCacheMap converts map into a binary string and loads the string into the database
// if the expires flag is false then the cached map will be added to the table with NULL
// for the expires tag and will not be removed when a CheckCache() occurs. If the expires
// flag is set to true a unix time will be written to the expires column of the record
// base on the configuration value of ExpireTime in minutes.
func (cache *Cache) SetCacheMap(cid string, data map[string]interface{}, expires bool) {
	var denc Data
	DebugHandler.Sys("SetCacheMap::"+cid, "Cache")
	gob.Register(Data{})
	cache.buf.Reset()
	denc.DataMap = data
	enc := gob.NewEncoder(&cache.buf)
	err := enc.Encode(&denc)
	if err != nil {
		LogHandler.Err(err, "Cache")
		DebugHandler.Err(err, "Cache", 3)
	}
	if !expires {
		cache.writeCache(cid, cache.buf.Bytes())
	} else {
		unixExpTime := time.Unix(0, time.Now().Add(time.Duration(cache.expireTime*time.Minute)).UnixNano())
		cache.writeCacheExp(cid, cache.buf.Bytes(), unixExpTime.Unix())
	}
}

// GetCacheMap returns requested cache map
func (cache *Cache) GetCacheMap(cid string) (bool, map[string]interface{}) {
	var data Data
	DebugHandler.Sys("GetCacheMap::"+cid, "Cache")
	gob.Register(Data{})
	exists := false
	cache.buf.Reset()
	dec := gob.NewDecoder(&cache.buf)
	cacheExists, cacheData := cache.readCache(cid)
	if cacheExists {
		_, err := cache.buf.Write(cacheData)
		if err != nil {
			LogHandler.Err(err, "Cache")
			DebugHandler.Err(err, "Cache", 3)
		}
		err = dec.Decode(&data)
		if err != nil {
			LogHandler.Err(err, "Cache")
			DebugHandler.Err(err, "Cache", 3)
		} else {
			exists = true
		}
	}
	return exists, data.DataMap
}

// WriteCache creates a new cache entry
func (cache *Cache) writeCache(cid string, data []byte) {
	//debug.DebugHandler.Sys("WriteCache", thisModule)
	_, err := cache.db.sql.Exec(`INSERT INTO cache (cid, data) VALUES (?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?`, cid, data, data)
	if err != nil {
		LogHandler.Err(err, "Cache")
		DebugHandler.Err(err, "Cache", 3)
	}
}

// WriteCacheExp creates a new cache entry that expires
func (cache *Cache) writeCacheExp(cid string, data []byte, expires int64) {
	//debug.DebugHandler.Sys("WriteCacheExp", thisModule)
	_, err := cache.db.sql.Exec(`INSERT INTO cache (cid, data, expires) VALUES (?, ?, ?)`+
		`ON DUPLICATE KEY UPDATE data = ?, expires = ?`, cid, data, expires, data, expires)
	if err != nil {
		LogHandler.Err(err, "Cache")
		DebugHandler.Err(err, "Cache", 3)
	}
}

// ReadCache returns a cache entry
func (cache *Cache) readCache(cid string) (bool, []byte) {
	var data []byte
	var expires sql.NullInt64
	var exists bool
	DebugHandler.Sys("ReadCache", "Cache")
	err := cache.db.sql.QueryRow(`SELECT data, expires FROM cache WHERE cid = ?`, cid).Scan(&data, &expires)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("CacheNotFound", "Cache")
		exists = false
	case err != nil:
		LogHandler.Err(err, "Cache")
		DebugHandler.Err(err, "Cache", 3)
		exists = false
	case expires.Valid:
		if expires.Int64 < time.Now().Unix() {
			exists = false
		} else {
			exists = true
		}
	default:
		exists = true
	}
	return exists, data
}

// ClearCache clears all records in cache table
func (cache *Cache) clearCache() {
	_, err := cache.db.sql.Exec(`TRUNCATE TABLE cache`)
	if err != nil {
		LogHandler.Err(err, "Cache")
		DebugHandler.Err(err, "Cache", 3)
	}
}

// CheckCache iterates through cache records checks if the record is expired
// if the record is expired the record is deleted, if the record has NULL for
// a value it is permanent and is skipped, if the record has not expired it is skipped
func (cache *Cache) checkCache() {
	var cid string
	var expires sql.NullInt64
	rows, err := cache.db.sql.Query(`SELECT cid, expires FROM cache`)
	switch {
	case err == sql.ErrNoRows:
		DebugHandler.Sys("ExpiredRecordsNotFound", "Cache")
	case err != nil && err.Error() != "EOF":
		DebugHandler.Err(err, "Cache", 3)
		LogHandler.Err(err, "Cache")
	default:
		for rows.Next() {
			rows.Scan(&cid, &expires)
			if expires.Valid && expires.Int64 < time.Now().Unix() {
				DebugHandler.Sys("RemovingExpired::"+cid, "Cache")
				_, err := cache.db.sql.Exec(`DELETE FROM cache WHERE cid = ?`, cid)
				if err != nil {
					DebugHandler.Err(err, "Cache", 3)
					LogHandler.Err(err, "Cache")
				}
			}
		}
	}
}
