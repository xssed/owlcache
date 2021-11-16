Example
```shell
package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/xssed/owlcache/counter"
)

//创建一个全局错误请求控制计数器
var MyCounter *counter.Counter

func main() {
	//初始化错误请求控制计数器
	MyCounter = counter.NewCounter()

	var do bool = true
	ticker := time.NewTicker(time.Second * time.Duration(20))
	go func() {
		for _ = range ticker.C {
			if do {
				do = false
			} else {
				do = true
			}
		}
	}()

	mcrts_exptime, _ := time.ParseDuration("5s") //假定睡眠时间，6秒   使用时-1
	mrmen_maxnum, _ := strconv.Atoi("2")         //假定最大请求数，2次   使用时-1
	//计数器假想，每一秒钟执行一次
	for i := 0; i < 1000000; i++ {
		time.Sleep(time.Second * time.Duration(1)) //人为阻塞

		key := "test_string"

		k := MyCounter.Exe(key, int64(mrmen_maxnum-1), mcrts_exptime)
		if k > 0 {
			fmt.Println(do, "允许执行", k)

			//do someing...

			if do {
				//执行成功-1
				MyCounter.Dec(key)
			}
		} else {
			fmt.Println(do, "不允许执行,睡眠状态", k)
		}

	}

}
```