package core

import (
	"testing"
)

func TestCache_InitCache(t *testing.T) {
	InitDebug(false)
	InitCron()
	InitDatabase()
	CacheHandler.CheckCache()
	CacheHandler.SetCacheMap("Test", GetOptions("Core"), true)
	exists, testCache := CacheHandler.GetCacheMap("Test")
	if exists {
		if testCache["Name"].(string) == "Core" {
			// Test Sucsessful
		} else {
			t.Error("CacheFailed")
		}
	} else {
		t.Error("CacheFailed")
	}
	CacheHandler.ClearAllCache()
}
