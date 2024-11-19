package cacheServer

import (
	"ProjectDemo/memoryCacheDemo/myCache"
	"time"
)

type MemCacheServer struct {
	cache *myCache.MemCache
}

func NewMemCache() *MemCacheServer {
	return &MemCacheServer{
		cache: myCache.NewMemCache(),
	}
}

// size: 1KB, 100KB, 1MB, 2MB, 1GB
func (cs *MemCacheServer) SetMaxMemory(size string) bool {
	return cs.cache.SetMaxMemory(size)
}

// 将value写入缓存
func (cs *MemCacheServer) Set(key string, val interface{}, expire ...time.Duration) bool {
	expireTs := time.Second * 0
	if len(expire) > 0 {
		expireTs = expire[0]
	}
	return cs.cache.Set(key, val, expireTs)
}

// 根据key值获取value
func (cs *MemCacheServer) Get(key string) (interface{}, bool) {
	return cs.cache.Get(key)
}

// 删除key值
func (cs *MemCacheServer) Del(key string) bool {
	return cs.cache.Del(key)
}

// 判断key是否存在
func (cs *MemCacheServer) Exists(key string) bool {
	return cs.cache.Exists(key)
}

// 清空所有key
func (cs *MemCacheServer) Flush() bool {
	return cs.cache.Flush()
}

// 获取缓存中所有key的数量
func (cs *MemCacheServer) keys() int64 {
	return cs.cache.Keys()
}
