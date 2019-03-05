package cache

//底层封装了sync.Map来实现线程安全的Map
//但是sync.Map中缺少很多方法，例如获取Map.dirty的方法等等

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

//创建一个缓存库模型
type BaseCache struct {
	CacheName    string
	KvStoreItems sync.Map
}

//增加一条内容
func (baseCache *BaseCache) Set(key interface{}, value interface{}, lifeTime time.Duration) bool {
	baseCache.KvStoreItems.Store(key, newKvStore(key, value, lifeTime))
	return true
}

//给一个key设置(或更新)过期
//返回否 代表这个key的内容不存在  正 代表成功
func (baseCache *BaseCache) Expire(key interface{}, lifeTime time.Duration) bool {
	v, ok := baseCache.KvStoreItems.Load(key)
	//如果存在
	if ok {
		kvstore := &KvStore{v.(*KvStore).Key, v.(*KvStore).Value, lifeTime, time.Now()} //v.(*KvStore).CreateTime 最早设计时想了很久。还是根据更新过期时间为起始时间计算过期吧
		baseCache.KvStoreItems.Store(key, kvstore)
		return true
	}
	return ok
}

//删除一条内容
func (baseCache *BaseCache) Delete(key interface{}) bool {
	baseCache.KvStoreItems.Delete(key)
	return true
}

//获取一条内容
func (baseCache *BaseCache) Get(key interface{}) (interface{}, bool) {
	v, ok := baseCache.KvStoreItems.Load(key)
	//如果存在
	if ok {
		//检查是否过期
		if !v.(*KvStore).IsExpired() {
			return v.(*KvStore).Value, ok
		} else {
			//过期key删除
			baseCache.Delete(key)
			return nil, false
		}
	}
	return nil, ok
}

//遍历集合
func (baseCache *BaseCache) GetKvStoreSlice() []*KvStore {

	var items []*KvStore

	baseCache.KvStoreItems.Range(func(k, v interface{}) bool {
		items = append(items, &KvStore{
			v.(*KvStore).Key,
			v.(*KvStore).Value,
			v.(*KvStore).LifeTime,
			v.(*KvStore).CreateTime})
		return true
	})

	return items

}

//判断一个key是否存在
func (baseCache *BaseCache) Exists(key interface{}) bool {
	_, ok := baseCache.Get(key)
	return ok
}

//刷新服务器上保存的内容 清空缓存
func (baseCache *BaseCache) Flush() bool {

	baseCache.KvStoreItems.Range(func(k, v interface{}) bool {
		//删除
		baseCache.Delete(k)
		return true
	})

	return true
}

//清除过期数据
func (baseCache *BaseCache) ClearExpireData() bool {

	baseCache.KvStoreItems.Range(func(k, v interface{}) bool {
		baseCache.Get(k)
		return true
	})

	return true
}

//返回缓存中数据项的数量
func (baseCache *BaseCache) Count() int {

	count := 0
	baseCache.KvStoreItems.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count

}

//将内存数据序列化到文件
func (baseCache *BaseCache) SaveToFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err = baseCache.SaveIoWriter(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

//将缓存数据写入io.Writer中
func (baseCache *BaseCache) SaveIoWriter(w io.Writer) (err error) {
	var mu sync.RWMutex

	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error Registering types with gob library")
		}
	}()

	mu.RLock()
	defer mu.RUnlock()

	var items []*KvStore

	items = baseCache.GetKvStoreSlice()
	//	for _, v := range items {
	//		gob.Register(&v)
	//		fmt.Println(v)
	//	}
	err = enc.Encode(&items)

	return err
}

//从文件中读取序列化的数据
func (baseCache *BaseCache) LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	if err = baseCache.LoadIoReader(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

//从io.Reader读取
func (baseCache *BaseCache) LoadIoReader(r io.Reader) error {
	var mu sync.RWMutex

	dec := gob.NewDecoder(r)
	var items []*KvStore
	err := dec.Decode(&items)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	//重新装载数据
	for _, v := range items {

		//fmt.Println(index, v)
		key := v.Key
		value := v.Value
		lifeTime := v.LifeTime
		createTime := v.CreateTime

		//检查是否过期
		if !v.IsExpired() {
			//没过期
			baseCache.KvStoreItems.Store(key, &KvStore{
				key,
				value,
				lifeTime,
				createTime})
		}

	}
	return nil
}
