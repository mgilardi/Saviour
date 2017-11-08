package cache

import (
  "encoding/gob"
  "config"
  "time"
  "modules/logger"
  "modules/database"
  "os"
  "bytes"
)

const (
  thisModule = "Cache"
)

type Cache struct {
  expireTime time.Duration
  db *database.Database
  options *config.Setting
  logger *logger.LogData
  buf bytes.Buffer
  encode *gob.Encoder
  decode *gob.Decoder
  memCache map[string]interface{}
}

func InitCache(conf *[]config.Setting, datab *database.Database, log *logger.LogData) *Cache {
  var cache Cache
  var err error
  cache.encode = gob.NewEncoder(&cache.buf)
  cache.decode = gob.NewDecoder(&cache.buf)
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

func (cache Cache) CheckCache() {
  //test
  os.Exit(0)
}
