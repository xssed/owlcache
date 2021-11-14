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
	ticker := time.NewTicker(time.Second * time.Duration(10))
	go func() {
		for _ = range ticker.C {
			if do {
				do = false
			} else {
				do = true
			}
		}
	}()

	mcrts_exptime, _ := time.ParseDuration("1s") //假定睡眠时间，2秒   使用时-1
	mrmen_maxnum, _ := strconv.Atoi("5")         //假定最大请求数，2次   使用时-1
	//计数器假想，每一秒钟执行一次
	for i := 0; i < 1000000; i++ {
		time.Sleep(time.Second * time.Duration(1)) //人为阻塞
		if !do {
			fmt.Println(MyCounter.Add("test_string", int64(mrmen_maxnum-1), mcrts_exptime))
		} else {
			fmt.Println("SUCCESS")
		}
	}

}
```