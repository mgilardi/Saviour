package core

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"strconv"
	"time"
)

// CacheHandler Global Access
var CacheHandler *Cache

// CacheObj interface for any struct containing the Cache method
type CacheObj interface {
	Cache() (string, map[string]interface{})
}

// CacheData binary encoding struct for the cache map
type CacheData struct {
	DataMap map[string]interface{}
}

// Cache struct holds the
type Cache struct {
	expireTime time.Duration
	db         *Database
}

// InitCache sets up the CacheHandler global variable
func InitCache(db *Database) {
	var cache Cache
	cache.db = db
	options := GetOptions("core")
	cache.expireTime = time.Duration(int(options["ExpireTime"].(float64)))
	CacheHandler = &cache
	CronHandler.Add(func() {
		CacheHandler.CheckCache()
	})
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
	cid, _ = obj.Cache()
	err = cache.db.sql.QueryRow(`SELECT data, expires FROM cache WHERE cid = ?`, cid).Scan(&dbData, &expires)
	switch {
	case err == sql.ErrNoRows:
		// Cache Row Not found
		Sys("CacheEntryNotFound::Creating", "Cache")
		cache.SetCache(obj)
		cache.db.sql.QueryRow(`SELECT data, expires FROM cache WHERE cid = ?`, cid).Scan(&dbData, &expires)
		exists = true
	case err != nil:
		Error(err, "Cache")
	case expires.Valid:
		if !(expires.Int64 < time.Now().Unix()) {
			exists = true
		} else {
			cache.SetCache(obj)
			cache.db.sql.QueryRow(`SELECT data, expires FROM cache WHERE cid = ?`, cid).Scan(&dbData, &expires)
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
			Error(err, "Cache")
		}
		err := decoder.Decode(&cacheData)
		if err != nil {
			Error(err, "Cache")
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
		Error(err, "Cache")
	}
	expTime = time.Unix(0, time.Now().Add(time.Duration(cache.expireTime*time.Minute)).UnixNano()).Unix()
	Sys("WritingCache::"+cid+":"+strconv.FormatInt(expTime, 10), "Cache")
	_, err = cache.db.sql.Exec(`INSERT INTO cache (cid, data, expires) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE data = ?, expires = ?`, cid, buf.Bytes(), expTime, buf.Bytes(), expTime)
	if err != nil {
		Error(err, "Cache")
	}
}

// ClearCache clears all records in cache table
func (cache *Cache) ClearCache() {
	_, err := cache.db.sql.Exec(`TRUNCATE TABLE cache`)
	if err != nil {
		Error(err, "Cache")
	}
	cache.db.ResetIncrement("cache")
}

// CheckCache iterates through cache records checks if the record is expired
// if the record is expired the record is deleted, if the record has NULL for
// a value it is permanent and is skipped, if the record has not expired it is skipped
func (cache *Cache) CheckCache() {
	var cid string
	var expires sql.NullInt64
	rows, err := cache.db.sql.Query(`SELECT cid, expires FROM cache`)
	Sys("CheckingForExpiredEntrys", "Cache")
	switch {
	case err == sql.ErrNoRows:
		Sys("ExpiredRecordsNotFound", "Cache")
	case err != nil && err.Error() != "EOF":
		Error(err, "Cache")
	default:
		for rows.Next() {
			rows.Scan(&cid, &expires)
			Sys("CheckingEntry::"+cid, "Cache")
			Sys("ExpireTime::"+strconv.FormatInt(expires.Int64, 10)+"::TimeNow::"+strconv.FormatInt(time.Now().Unix(), 10), "Cache")
			if expires.Valid && expires.Int64 < time.Now().Unix() {
				Sys("RemovingExpired::"+cid, "Cache")
				_, err := cache.db.sql.Exec(`DELETE FROM cache WHERE cid = ?`, cid)
				if err != nil {
					Error(err, "Cache")
				}
			}
		}
		cache.db.ResetIncrement("cache")
	}
}
