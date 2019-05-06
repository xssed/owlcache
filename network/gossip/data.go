package gossip

import (
	"github.com/xssed/owlcache/cache"
	"github.com/xssed/owlcache/network"
)

//Key/Value数据结构
type Data struct{}

//创建
func NewData() *Data {
	return &Data{}
}

//设置Key数据
func (d *Data) Set(key string, value string, expires time.Duration) bool {
	ok := network.BaseCacheDB.Set(key, value, expires) //Expires time.Duration //int64
	if ok {
		return true
	} else {
		return false
	}
}

//获取Key数据
func (d *Data) Get(key string) string {
	var str string
	if v, found := network.BaseCacheDB.GetKvStore(key); found {
		str = v.(*cache.KvStore).Value
	} else {
		str = ""
	}
	return str
}

//为Key设置过期时间
func (d *Data) Expire(key string, expires time.Duration) bool {
	ok := network.BaseCacheDB.Expire(key, expires) //Expires time.Duration //int64
	if ok {
		return true
	} else {
		return false
	}
}

//删除Key数据
func (d *Data) Delete(key string) bool {
	ok := network.BaseCacheDB.Delete(key)
	if ok {
		return true
	} else {
		return false
	}
}

//获取整个数据集合
func (d *Data) GetItems() []*cache.KvStore {
	m := network.BaseCacheDB.GetKvStoreSlice()
	return m
}
