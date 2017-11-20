package main

import (
	"bytes"
	"encoding/gob"
	"strings"
	"time"
)

const (
	thisModuleCache = "Cache"
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
	options := GetOptions(thisModuleCache)
	newCache.expireTime = time.Duration(int(options["ExpireTime"].(float64)))
	newCache.CacheOptions()
	CacheHandler = &newCache
	CronHandler.Add("CacheCheck", true, func() {
		CacheHandler.CheckCache()
		CacheHandler.CacheOptions()
	})
}

// CacheOptions loads the modules configuration files into cache
func (cache *Cache) CacheOptions() {
	DebugHandler.Sys("CacheOptions", thisModuleCache)
	allOptions := GetAllOptions()
	for _, opt := range allOptions {
		cache.SetCacheMap(strings.ToLower(opt["Name"].(string)+":config"), opt, false)
	}
}

// CheckCache checks for expired cache entrys
func (cache *Cache) CheckCache() {
	DebugHandler.Sys("CheckCacheForExpired", thisModuleCache)
	cache.db.CheckCache()
}

// ClearAllCache removes all the records from the cache and reloads module options
func (cache *Cache) ClearAllCache() {
	DebugHandler.Sys("ClearAllCache", thisModuleCache)
	cache.db.ClearCache()
	cache.CacheOptions()
}

// SetCacheMap converts map into a binary string and loads the string into the database
// if the expires flag is false then the cached map will be added to the table with NULL
// for the expires tag and will not be removed when a CheckCache() occurs. If the expires
// flag is set to true a unix time will be written to the expires column of the record
// base on the configuration value of ExpireTime in minutes.
func (cache *Cache) SetCacheMap(cid string, data map[string]interface{}, expires bool) {
	var denc Data
	DebugHandler.Sys("SetCacheMap::"+cid, thisModuleCache)
	gob.Register(Data{})
	cache.buf.Reset()
	denc.DataMap = data
	enc := gob.NewEncoder(&cache.buf)
	err := enc.Encode(&denc)
	if err != nil {
		LogHandler.Err(err, thisModuleCache)
		DebugHandler.Err(err, thisModuleCache, 3)
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
	DebugHandler.Sys("GetCacheMap::"+cid, thisModuleCache)
	gob.Register(Data{})
	exists := false
	cache.buf.Reset()
	dec := gob.NewDecoder(&cache.buf)
	cacheExists, cacheData := cache.db.ReadCache(cid)
	if cacheExists {
		_, err := cache.buf.Write(cacheData)
		if err != nil {
			LogHandler.Err(err, thisModuleCache)
			DebugHandler.Err(err, thisModuleCache, 3)
		}
		err = dec.Decode(&data)
		if err != nil {
			LogHandler.Err(err, thisModuleCache)
			DebugHandler.Err(err, thisModuleCache, 3)
		} else {
			exists = true
		}
	}
	return exists, data.DataMap
}
