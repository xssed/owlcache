package queue

//thread-safe queue
import (
	//"fmt"
	"sync"
)

type Element interface{}

type QueueI interface {
	Push(e Element) //向队列中添加元素
	Pop() Element   //取出队列中最前面的元素并且从队列中移除
	Clear() bool    //清空队列
	Size() int      //获取队列的元素个数
	IsEmpty() bool  //判断队列是否是空
}

type Queue struct {
	element []Element
	lock    sync.Mutex
}

func New() *Queue {
	return &Queue{}
}

//向队列中添加元素
func (q *Queue) Push(e Element) {
	q.lock.Lock()
	q.element = append(q.element, e)
	q.lock.Unlock()
}

//取出队列中最前面的元素并且从队列中移除
func (q *Queue) Pop() Element {

	if q.IsEmpty() {
		//fmt.Println("queue is empty!")
		return nil
	}
	q.lock.Lock()
	firstElement := q.element[0]
	q.element = q.element[1:]
	q.lock.Unlock()

	return firstElement
}

//清空队列
func (q *Queue) Clear() bool {

	if q.IsEmpty() {
		//fmt.Println("queue is empty!")
		return false
	}
	size := q.Size()
	q.lock.Lock()
	for i := 0; i < size; i++ {
		q.element[i] = nil
	}
	q.element = nil
	q.lock.Unlock()

	return true
}

//获取队列的元素个数
func (q *Queue) Size() int {

	q.lock.Lock()
	m := len(q.element)
	q.lock.Unlock()

	return m
}

//判断队列是否是空
func (q *Queue) IsEmpty() bool {

	q.lock.Lock()
	m := len(q.element)
	q.lock.Unlock()

	if m == 0 {
		return true
	}
	return false
}
