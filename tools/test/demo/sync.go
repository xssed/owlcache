//package main

//import (
//	"fmt"
//)

//func main() {
//	ch := make(chan struct{})
//	count := 2 // count 表示活动的协程个数
//	go func() {
//		fmt.Println("Goroutine 1")
//		ch <- struct{}{} // 协程结束，发出信号
//	}()
//	go func() {
//		fmt.Println("Goroutine 2")
//		ch <- struct{}{} // 协程结束，发出信号
//	}()
//	for range ch {
//		// 每次从ch中接收数据，表明一个活动的协程结束
//		count--
//		// 当所有活动的协程都结束时，关闭管道
//		if count == 0 {
//			close(ch)
//		}
//	}
//}
//package main

//import (
//	"fmt"
//	"sync"
//)

//func main() {
//	var wg sync.WaitGroup
//	wg.Add(2) // 因为有两个动作，所以增加2个计数
//	go func() {
//		fmt.Println("Goroutine 1")
//		wg.Done() // 操作完成，减少一个计数
//	}()
//	go func() {
//		fmt.Println("Goroutine 2")
//		wg.Done() // 操作完成，减少一个计数
//	}()
//	wg.Wait() // 等待，直到计数为0
//}
// package main

// import (
// 	"fmt"
// 	"time"
// )

// type Demo struct {
// 	input         chan string
// 	output        chan string
// 	max_goroutine chan int
// }

// func NewDemo() *Demo {
// 	d := new(Demo)
// 	d.input = make(chan string, 24)
// 	d.output = make(chan string, 24)
// 	d.max_goroutine = make(chan int, 20)
// 	return d
// }

// func (this *Demo) Goroutine() {
// 	var i = 1000
// 	for {
// 		this.input <- time.Now().Format("2006-01-02 15:04:05")
// 		time.Sleep(time.Second * 1)
// 		if i < 0 {
// 			break
// 		}
// 		i--
// 	}
// 	close(this.input)
// }

// func (this *Demo) Handle() {
// 	for t := range this.input {
// 		fmt.Println("datatime is :", t)
// 		this.output <- t
// 	}
// }

// func main() {
// 	demo := NewDemo()
// 	go demo.Goroutine()
// 	demo.Handle()
// }
package main

import (
	"fmt"
)

func main() {
	ch := make(chan struct{})
	count := 2 // count 表示活动的协程个数
	go func() {
		fmt.Println("Goroutine 1")
		ch <- struct{}{} // 协程结束，发出信号
	}()
	go func() {
		fmt.Println("Goroutine 2")
		ch <- struct{}{} // 协程结束，发出信号
	}()
	for range ch {
		// 每次从ch中接收数据，表明一个活动的协程结束
		count--

		fmt.Println(ch)

		// 当所有活动的协程都结束时，关闭管道
		if count == 0 {
			close(ch)
		}
	}

}
