package database

import (
	"bytes"
	"config"
	"encoding/gob"
	"modules/logger"
	"strings"
	"time"
)

type Data struct {
	DataMap map[string]interface{}
}

type Cache struct {
	expireTime time.Duration
	options    map[string]interface{}
	buf        bytes.Buffer
	db         *Database
}

func InitCache(db *Database) *Cache {
	var cache Cache
	cache.db = db
	cache.options = config.GetOptions(thisModule)
	cache.expireTime = time.Duration(int(cache.options["ExpireTime"].(float64)))
	cache.cacheOptions()
	return &cache
}

func (cache *Cache) cacheOptions() {
	allOptions := config.GetAllOptions()
	for _, opt := range allOptions {
		cache.SetCacheMap(strings.ToLower(opt["Name"].(string)+":config"), opt, false)
	}
}

func (cache *Cache) CheckCache() {
}

func (cache *Cache) SetCacheMap(cid string, data map[string]interface{}, expires bool) {
	var denc Data
	gob.Register(Data{})
	cache.buf.Reset()
	denc.DataMap = data
	enc := gob.NewEncoder(&cache.buf)
	err := enc.Encode(&denc)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
	}
	cache.db.WriteCache(cid, cache.buf.Bytes(), true)
}

func (cache *Cache) GetCacheMap(cid string) (bool, map[string]interface{}) {
	var data Data
	gob.Register(Data{})
	exists := true
	cache.buf.Reset()
	dec := gob.NewDecoder(&cache.buf)
	_, err := cache.buf.Write(cache.db.ReadCache(cid))
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
		exists = false
	}
	err = dec.Decode(&data)
	if err != nil {
		logger.Error(err.Error(), thisModule, 3)
		exists = false
	}
	return exists, data.DataMap
}
