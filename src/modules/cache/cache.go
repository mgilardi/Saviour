package cache

import (
  "encoding/gob"
  "config"
  "time"
  "modules/logger"
  "modules/database"
  "os"
)

const (
  thisModule = "Cache"
)

type dbObj struct {
  Object interface{}
  Created int64
  Expires int64
}

type memObj struct {
  Object interface{}
  Created int64
  Expires int64
}

type Cache struct {
  expireTime time.Duration
  db *database.Database
  options *config.Setting
  logger *logger.LogData
  encode gob.GobEncoder
  memCache map[string]memObj
}

func InitCache(conf *[]config.Setting, datab *database.Database, log *logger.LogData ) *Cache {
  var cache Cache
  var err error
  cache.logger = log
  cache.logger.SystemMessage("Starting", thisModule)
  cache.db = datab
  err, cache.options = config.GetSettingModule(thisModule, conf)
  if err != nil {
    cache.logger.Error(err.Error(), thisModule, 3)
    cache.logger.Error("CannotLoadSettings", thisModule, 1)
  }
  cache.expireTime = time.Duration(int(cache.options.FindValue("ExpireTime").(float64)))
  return &cache
}

func (cache Cache) Set(key string, data interface{}) {
  nObject := dbObj{data, time.Now().UnixNano(),
    (time.Now().UnixNano() + time.Duration(time.Minute * cache.expireTime).Nanoseconds())}
  gob.Register(nObject)
}

func (cache Cache) CheckCache() {
  //test
  os.Exit(0)
}
