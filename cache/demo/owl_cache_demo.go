package main

import (
	"fmt"
	"time"

	"github.com/xssed/owlcache/cache"
)

func main() {

	c := cache.NewCache("owlcache") //随意指定一个名字

	expiration, _ := time.ParseDuration("2s")
	key1 := "hello world!"
	c.Set("key1", key1, expiration) //两秒过期的缓存

	key2 := "测试啊!"
	c.Set("key2", key2, 0) //永不过期

	fmt.Println("共有", c.Count(), "个数据")

	if v, found := c.Get("key1"); found {
		fmt.Println("found key1:", v)
	} else {
		fmt.Println("not found key1")
	}

	if v2, found2 := c.Get("key2"); found2 {
		fmt.Println("found key2:", v2)
	} else {
		fmt.Println("not found key2")
	}

	s, _ := time.ParseDuration("4s")
	time.Sleep(s)

	if v3, found3 := c.Get("key1"); found3 {
		fmt.Println("found key1:", v3)
	} else {
		fmt.Println("not found key1")
	}

	if v4, found4 := c.Get("key2"); found4 {
		fmt.Println("found key2:", v4)
	} else {
		fmt.Println("not found key2")
	}

	fmt.Println("========================")
	c.Set("key3", "大吉大利", 0)               //永不过期
	c.Set("key4", "新年快乐", 200*time.Second) //200S后过期
	c.Set("key5", "项目早日完工", 0)             //永不过期
	c.Set("key6", "测试删除", 0)               //永不过期

	fmt.Println("共有", c.Count(), "个数据")

	c.Delete("key6")

	fmt.Println("共有", c.Count(), "个数据")
	fmt.Println("key3存在嘛？", c.Exists("key3"))
	fmt.Println("key4存在嘛？", c.Exists("key4"))
	fmt.Println("key6存在嘛？", c.Exists("key6"))
	fmt.Println("========================")

	err := c.SaveToFile("./owlcache.db")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("共有", c.Count(), "个数据")

	c.Flush() //清除所有

	fmt.Println("共有", c.Count(), "个数据")

	err2 := c.LoadFromFile("./owlcache.db")
	if err2 != nil {
		fmt.Println(err2)
	}

	fmt.Println("共有", c.Count(), "个数据")

	for index, v := range c.GetKvStoreSlice() {
		fmt.Println(index, v)
	}

}
