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
// package main

// import (
// 	"fmt"
// )

// func main() {
// 	ch := make(chan struct{})
// 	count := 2 // count 表示活动的协程个数
// 	go func() {
// 		fmt.Println("Goroutine 1")
// 		ch <- struct{}{} // 协程结束，发出信号
// 	}()
// 	go func() {
// 		fmt.Println("Goroutine 2")
// 		ch <- struct{}{} // 协程结束，发出信号
// 	}()
// 	for range ch {
// 		// 每次从ch中接收数据，表明一个活动的协程结束
// 		count--

// 		fmt.Println(ch)

// 		// 当所有活动的协程都结束时，关闭管道
// 		if count == 0 {
// 			close(ch)
// 		}
// 	}

// }
package main

import (
	"fmt"
	"sync"
)

type Set struct {
	m map[interface{}]bool
	sync.RWMutex
}

func New() *Set {
	return &Set{
		m: map[interface{}]bool{},
	}
}

func (s *Set) Add(item interface{}) {
	//写锁
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *Set) Remove(item interface{}) {
	//写锁
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *Set) Has(item interface{}) bool {
	//允许读
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *Set) List() []interface{} {
	//允许读
	s.RLock()
	defer s.RUnlock()
	var outList []interface{}
	for value := range s.m {
		outList = append(outList, value)
	}
	return outList
}

func (s *Set) Len() int {
	return len(s.List())
}

func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

func (s *Set) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func main() {

	s := New()
	wait := sync.WaitGroup{}

	go func() {
		wait.Add(1)
		defer wait.Done()
		s.Add(1)
		s.Add("2")
		s.Add("3")
	}()

	go func() {
		wait.Add(1)
		defer wait.Done()
		s.Add(1)
		s.Add("2")
		s.Add("3")
	}()

	go func() {
		wait.Add(1)
		defer wait.Done()
		s.Add(3.1415926)
		//fmt.Println(s.List())
		s.Remove("2")
		//fmt.Println(s.List())
	}()

	go func() {
		wait.Add(1)
		defer wait.Done()
		if s.Has("2") {
			fmt.Println("2 exist")
		} else {
			fmt.Println("2 not exist")
		}
	}()

	go func() {
		wait.Add(1)
		defer wait.Done()
		if s.Has(3.1415926) {
			fmt.Println("3.1415926 exist")
		} else {
			fmt.Println("3.1415926 not exist")
		}
	}()

	for i := 99; i < 999; i++ {
		s.Add(i)
	}

	fmt.Println("main gorotue :")
	fmt.Println("clear before ")
	fmt.Println("len == ", s.Len())
	fmt.Println("Is empty:", s.IsEmpty())
	s.Clear()
	fmt.Println("clear after")
	fmt.Println("is empty:", s.IsEmpty())
	fmt.Println("len == ", s.Len())
	wait.Wait()
}
