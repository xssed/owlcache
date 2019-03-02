package cache

import (
	"sync"
)

func NewCache(cacheName string) *BaseCache {

	value := &BaseCache{cacheName, sync.Map{}}
	return value

}
