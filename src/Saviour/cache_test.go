package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestInitCache(t *testing.T) {
	type args struct {
		db *Database
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitCache(tt.args.db)
		})
	}
}

func TestCache_cacheOptions(t *testing.T) {
	type fields struct {
		expireTime time.Duration
		options    map[string]interface{}
		buf        bytes.Buffer
		db         *Database
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				expireTime: tt.fields.expireTime,
				options:    tt.fields.options,
				buf:        tt.fields.buf,
				db:         tt.fields.db,
			}
			cache.cacheOptions()
		})
	}
}

func TestCache_CheckCache(t *testing.T) {
	type fields struct {
		expireTime time.Duration
		options    map[string]interface{}
		buf        bytes.Buffer
		db         *Database
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				expireTime: tt.fields.expireTime,
				options:    tt.fields.options,
				buf:        tt.fields.buf,
				db:         tt.fields.db,
			}
			cache.CheckCache()
		})
	}
}

func TestCache_ClearAllCache(t *testing.T) {
	type fields struct {
		expireTime time.Duration
		options    map[string]interface{}
		buf        bytes.Buffer
		db         *Database
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				expireTime: tt.fields.expireTime,
				options:    tt.fields.options,
				buf:        tt.fields.buf,
				db:         tt.fields.db,
			}
			cache.ClearAllCache()
		})
	}
}

func TestCache_SetCacheMap(t *testing.T) {
	type fields struct {
		expireTime time.Duration
		options    map[string]interface{}
		buf        bytes.Buffer
		db         *Database
	}
	type args struct {
		cid     string
		data    map[string]interface{}
		expires bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				expireTime: tt.fields.expireTime,
				options:    tt.fields.options,
				buf:        tt.fields.buf,
				db:         tt.fields.db,
			}
			cache.SetCacheMap(tt.args.cid, tt.args.data, tt.args.expires)
		})
	}
}

func TestCache_GetCacheMap(t *testing.T) {
	type fields struct {
		expireTime time.Duration
		options    map[string]interface{}
		buf        bytes.Buffer
		db         *Database
	}
	type args struct {
		cid string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  map[string]interface{}
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &Cache{
				expireTime: tt.fields.expireTime,
				options:    tt.fields.options,
				buf:        tt.fields.buf,
				db:         tt.fields.db,
			}
			got, got1 := cache.GetCacheMap(tt.args.cid)
			if got != tt.want {
				t.Errorf("Cache.GetCacheMap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Cache.GetCacheMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
