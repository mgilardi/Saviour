package main

import (
	"testing"
)

func TestCache_InitCache(t *testing.T) {
	InitDebug(false)
	InitCron()
	InitDatabase()
	CacheHandler.CacheOptions()
	CacheHandler.CheckCache()
	CacheHandler.SetCacheMap("Test", GetOptions("Cache"), true)
	exists, testCache := CacheHandler.GetCacheMap("Test")
	if exists {
		if testCache["Name"].(string) == "Cache" {
			// Test Sucsessful
		} else {
			t.Error("CacheFailed")
		}
	} else {
		t.Error("CacheFailed")
	}
	CacheHandler.ClearAllCache()
}
