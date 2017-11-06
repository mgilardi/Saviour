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

func (cache Cache) Set(key string, data interface{}) {
  err := cache.encode.Encode(data)
  if err != nil {
    //
  }
  cache.db.WriteCache(key, cache.buf.Bytes(), time.Now().UnixNano(), time.Duration(time.Minute * cache.expireTime).Nanoseconds())
  cache.logger.SystemMessage("Wrote::" + key, thisModule)
  cache.buf.Reset()
}

func (cache Cache) Get(key string) interface{} {
  err, data := cache.db.ReadCache(key)
  err = cache.decode.Decode(data)
  if err != nil {
    //
  }
  newdata := cache.buf
  cache.buf.Reset()
  return newdata
}

func (cache Cache) CheckCache() {
  //test
  os.Exit(0)
}
