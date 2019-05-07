thread-safe Queue Example
```shell

package main

import (
	"fmt"
	"strconv"

	"github.com/xssed/owlcache/queue"
)

func main() {
	queue := queue.New()

	for i := 0; i < 50; i++ {
		queue.Push(strconv.Itoa(i) + "测试")
	}
	queue.Push("51测试")
	fmt.Println("元素个数:", queue.Size())
	fmt.Println("移除最前面的元素：", queue.Pop())
	fmt.Println("移除最前面的元素：", queue.Pop())
	queue.Push("52测试")
	fmt.Println("移除最前面的元素：", queue.Pop())
	fmt.Println("移除最前面的元素：", queue.Pop())
	fmt.Println("元素个数:", queue.Size())
	fmt.Println("清空队列：", queue.Clear())
	fmt.Println("元素个数:", queue.Size())

}


```