## USE  

```shell
package main

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

var mc *memcache.Client

func main() {

	mc = memcache.New("127.0.0.1:11211", "192.168.0.50:11211")

	Set("hello", "world")
	key1, err := Get("hello")
	if err == nil {
		fmt.Println("hello:", string(key1.Value))
	}

	Add("hello1", "world1")
	key2, err1 := Get("hello1")
	if err1 == nil {
		fmt.Println("hello1:", string(key2.Value))
	}

	Replace("hello", "new world")
	key3, err2 := Get("hello")
	if err2 == nil {
		fmt.Println("hello:", string(key3.Value))
	}

	ok := Delete("hello1")
	if ok {
		fmt.Println("Delete hello1 OK")
	} else {
		fmt.Println("Delete hello1 error")
	}

}

//Set方法
func Set(key, value string) bool {
	err := mc.Set(&memcache.Item{Key: key, Value: []byte(value)})
	if err != nil {
		fmt.Println("Set() error:", err)
		return false
	}
	return true
}

//Get方法
func Get(key string) (*memcache.Item, error) {
	it, err := mc.Get(key)
	if err != nil {
		fmt.Println("Get() error:", err)
		return nil, err
	}
	return it, nil
}

//Add方法
func Add(key, value string) bool {
	err := mc.Add(&memcache.Item{Key: key, Value: []byte(value)})
	if err != nil {
		fmt.Println("Add() error:", err)
		return false
	}
	return true
}

//Replace方法
func Replace(key, value string) bool {
	err := mc.Replace(&memcache.Item{Key: key, Value: []byte(value)})
	if err != nil {
		fmt.Println("Replace() error:", err)
		return false
	}
	return true
}

//Delete方法
func Delete(key string) bool {
	err := mc.Delete(key)
	if err != nil {
		fmt.Println("Delete() error:", err)
		return false
	}
	return true
}

```