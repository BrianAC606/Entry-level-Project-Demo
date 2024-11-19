package myCache

import (
	"sync"
	"time"
)

type MemCache struct {
	//最大缓存大小
	maxMemoryCacheSize int64
	//当前缓存大小
	currentMemoryCacheSize int64
	//最大缓存大小
	maxMemoryCacheSizeStr string
	//用于线程同步的读写锁
	locker sync.RWMutex
	//缓存内容
	values map[string]*memCacheValue
	//清理过期key值时间间隔
	clearExpireKeyTime time.Duration
}

type memCacheValue struct {
	size int64
	//过期时间
	expireTime time.Time
	//生存时间
	lastAccessTime time.Duration
	//缓存值
	val interface{}
}

func NewMemCache() (mem *MemCache) {
	defer func() {
		go mem.delExpireKey()
	}()
	return &MemCache{
		maxMemoryCacheSizeStr:  "",
		maxMemoryCacheSize:     0,
		values:                 make(map[string]*memCacheValue),
		currentMemoryCacheSize: 0,
		locker:                 sync.RWMutex{},
		clearExpireKeyTime:     time.Second * 10,
	}
}

// size: 1KB, 100KB, 1MB, 2MB, 1GB
func (mem *MemCache) SetMaxMemory(size string) bool {
	mem.maxMemoryCacheSize, mem.maxMemoryCacheSizeStr = parseSize(size)
	println("SetMaxMemory:", mem.maxMemoryCacheSize, "string:", mem.maxMemoryCacheSizeStr)
	return true
}

// 将value写入缓存
func (mem *MemCache) Set(key string, val interface{}, expire time.Duration) (is bool) {
	v := &memCacheValue{
		size:           CalSize(val),
		val:            val,
		expireTime:     time.Now().Add(expire),
		lastAccessTime: expire,
	}
	if mem.Exists(key) {
		mem.Del(key)
	}
	mem.locker.Lock()
	mem.values[key] = v
	mem.currentMemoryCacheSize += v.size
	mem.locker.Unlock()
	if mem.currentMemoryCacheSize >= mem.maxMemoryCacheSize {
		mem.Del(key)
		panic("max memory cache size exceeded")
	}
	//defer func() {
	//	if err := recover(); err != nil {
	//		fmt.Println(err)
	//		is = false
	//	}
	//}()
	return true
}

// 根据key值获取value
func (mem *MemCache) Get(key string) (interface{}, bool) {
	mem.locker.RLock()
	defer mem.locker.RUnlock()
	value, ok := mem.values[key]
	if ok {
		if value.lastAccessTime > 0 && value.expireTime.Before(time.Now()) {
			mem.Del(key)
			return nil, false
		}
	}
	println("Get:", key)
	return value, ok
}

// 删除key值
func (mem *MemCache) Del(key string) bool {
	if mem.Exists(key) {
		mem.locker.Lock()
		defer mem.locker.Unlock()
		mem.currentMemoryCacheSize -= mem.values[key].size
		delete(mem.values, key)
	}
	return true
}

// 判断key是否存在
func (mem *MemCache) Exists(key string) bool {
	mem.locker.RLock()
	defer mem.locker.RUnlock()
	_, ok := mem.values[key]
	return ok
}

// 清空所有key
func (mem *MemCache) Flush() bool {
	mem.locker.Lock()
	defer mem.locker.Unlock()
	//比较低效的删除方式
	//for k, _ := range mem.values {
	//	delete(mem.values, k)
	//}
	mem.values = make(map[string]*memCacheValue, 0)
	mem.currentMemoryCacheSize = 0
	mem.maxMemoryCacheSize = 0
	mem.maxMemoryCacheSizeStr = ""
	return true
}

// 获取缓存中所有key的数量
func (mem *MemCache) Keys() int64 {
	mem.locker.RLock()
	defer mem.locker.RUnlock()
	return int64(len(mem.values))
}

func (mem *MemCache) delExpireKey() {
	timeTicker := time.NewTicker(mem.clearExpireKeyTime)
	defer timeTicker.Stop()
	for {
		select {
		case <-timeTicker.C:
			for key, value := range mem.values {
				if value.lastAccessTime > 0 && time.Now().Sub(value.expireTime) > mem.clearExpireKeyTime {
					mem.Del(key)
				}
			}
		}
	}
}
