package database

import (
	"bytes"
	"config"
	"encoding/gob"
	"modules/debug"
	"strings"
	"time"
)

// DBCache global cache access
var DBCache *Cache

// Data holds the data input/output of the binary encoder/decoder
type Data struct {
	DataMap map[string]interface{}
}

// Cache holds the options, buffer, and db access
type Cache struct {
	expireTime time.Duration
	options    map[string]interface{}
	buf        bytes.Buffer
	db         *Database
}

// InitCache constructs the cache object
func InitCache(db *Database) {
	var newCache Cache
	debug.Dbg.Sys("Starting", thisModule)
	newCache.db = db
	newCache.options = config.GetOptions(thisModule)
	newCache.expireTime = time.Duration(int(newCache.options["ExpireTime"].(float64)))
	newCache.cacheOptions()
	DBCache = &newCache
}

// cacheOptions loads the modules configuration files into cache
func (cache *Cache) cacheOptions() {
	debug.Dbg.Sys("CacheOptions", thisModule)
	allOptions := config.GetAllOptions()
	for _, opt := range allOptions {
		cache.SetCacheMap(strings.ToLower(opt["Name"].(string)+":config"), opt, false)
	}
}

// CheckCache checks for expired cache entrys
func (cache *Cache) CheckCache() {
	debug.Dbg.Sys("CheckCacheForExpired", thisModule)
	cache.db.CheckCache()
}

// ClearAllCache removes all the records from the cache and reloads module options
func (cache *Cache) ClearAllCache() {
	debug.Dbg.Sys("ClearAllCache", thisModule)
	cache.db.ClearCache()
	cache.cacheOptions()
}

// SetCacheMap converts map into a binary string and loads the string into the database
// if the expires flag is false then the cached map will be added to the table with NULL
// for the expires tag and will not be removed when a CheckCache() occurs. If the expires
// flag is set to true a unix time will be written to the expires column of the record
// base on the configuration value of ExpireTime in minutes.
func (cache *Cache) SetCacheMap(cid string, data map[string]interface{}, expires bool) {
	var denc Data
	debug.Dbg.Sys("SetCacheMap::"+cid, thisModule)
	gob.Register(Data{})
	cache.buf.Reset()
	denc.DataMap = data
	enc := gob.NewEncoder(&cache.buf)
	err := enc.Encode(&denc)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
	}
	if !expires {
		cache.db.WriteCache(cid, cache.buf.Bytes())
	} else {
		unixExpTime := time.Unix(0, time.Now().Add(time.Duration(cache.expireTime*time.Minute)).UnixNano())
		cache.db.WriteCacheExp(cid, cache.buf.Bytes(), unixExpTime.Unix())
	}
}

// GetCacheMap returns requested cache map
func (cache *Cache) GetCacheMap(cid string) (bool, map[string]interface{}) {
	var data Data
	debug.Dbg.Sys("GetCacheMap::"+cid, thisModule)
	cache.CheckCache()
	gob.Register(Data{})
	exists := true
	cache.buf.Reset()
	dec := gob.NewDecoder(&cache.buf)
	_, err := cache.buf.Write(cache.db.ReadCache(cid))
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
		exists = false
	}
	err = dec.Decode(&data)
	if err != nil {
		LogDB.Err(err, thisModule)
		debug.Dbg.Err(err, thisModule, 3)
		exists = false
	}
	return exists, data.DataMap
}
