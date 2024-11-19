package main

import "ProjectDemo/memoryCacheDemo/myCache"

func main() {
	cache := myCache.NewMemCache()
	cache.SetMaxMemory("100MB")
	cache.Set("int", 1, 10)
	cache.Set("bool", false, 10)
	cache.Set("data", map[string]interface{}{"a": 1}, 10)
	cache.Get("int")
	cache.Del("int")
	cache.Flush()
	cache.Keys()
}
