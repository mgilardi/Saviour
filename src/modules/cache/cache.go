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
  options map[string]interface{}
  buf bytes.Buffer
  encode *gob.Encoder
  decode *gob.Decoder
  memCache map[string]interface{}
}

func InitCache(datab *database.Database) *Cache {
  var cache Cache
  var err error
  cache.encode = gob.NewEncoder(&cache.buf)
  cache.decode = gob.NewDecoder(&cache.buf)
  logger.SystemMessage("Starting", thisModule)
  cache.db = datab
  cache.options = config.GetOptions(thisModule)
  if err != nil {
    logger.Error(err.Error(), thisModule, 3)
    logger.Error("CannotLoadSettings", thisModule, 1)
  }
  cache.expireTime = time.Duration(int(cache.options["ExpireTime"].(float64)))
  return &cache
}

func (cache Cache) CheckCache() {
  //test
  os.Exit(0)
}
