package lib

import (
	"sync"
	"time"
)
//缓存节点
type cacheNode struct {
	Data       interface{}
	Expiration time.Time
}
//全局缓存
var GlobalCacheMap map[string]cacheNode
var globalCacheLock sync.RWMutex

func init() {
	GlobalCacheMap = make(map[string]cacheNode)
	NewCache().Clean()
}

type Cache struct {
}

func NewCache() *Cache {
	return &Cache{}
}
//定期清理过期缓存
func (this *Cache) Clean() {
	go func() {
		defer CatchPanic()
		for {
			globalCacheLock.Lock()
			if len(GlobalCacheMap) > 0 {
				for key, node := range GlobalCacheMap {
					if time.Now().After(node.Expiration) {
						delete(GlobalCacheMap, key)
					}
				}
			}
			globalCacheLock.Unlock()
			time.Sleep(30 * time.Second)
		}
	}()
}
//放置缓存 value值为interface
func (this *Cache) Set(key string, value interface{}, expiration time.Duration) {
	globalCacheLock.Lock()
	defer globalCacheLock.Unlock()
	GlobalCacheMap[key] = cacheNode{Data: value, Expiration: time.Now().Add(expiration)}

}

func (this *Cache) Get(key string) interface{} {
	globalCacheLock.RLock()
	v, ok := GlobalCacheMap[key]
	globalCacheLock.RUnlock()
	if !ok {
		return nil
	}

	if time.Now().Before(v.Expiration) {
		return v.Data
	}
	return nil
}
//获取放置时为map的缓存
func (this *Cache) GetMap(key string) map[string]string {
	v := this.Get(key)
	if v == nil {
		return nil
	}
	retval, ok := v.(map[string]string)
	if ok {
		return retval
	}
	return nil
}
