package core

import (
	"fmt"
	"testing"
)

func TestCache_InitCache(t *testing.T) {
	initTest()
	exists, testCache := CacheHandler.GetCache(testHandler)
	if !exists {
		t.Error("GetCacheFailed")
	}
	fmt.Println("Output:" + testCache["test"].(string))
	CacheHandler.CheckCache()
	CacheHandler.ClearCache()
}
