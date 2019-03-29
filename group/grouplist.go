//这里是一个将Slice切片封装成List结构体的轮子
//为什么没有使用官方包的"container/list"？主要是测试后发现性能的差异，所以应用场景的使用上也要有所差异
//声明：测试结果为个人电脑的测试结果，仅供参考。
//===========================================================
//t := time.Now()
//sli := make([]int, 10)
//for i := 0; i < 1*100000*1000; i++ {
//	sli = append(sli, 1)
//}
//fmt.Println("slice 创建速度：" + time.Now().Sub(t).String())

//t = time.Now()
//l := list.New()
//for i := 0; i < 1*100000*1000; i++ {
//	l.PushBack(1)
//}
//fmt.Println("list 创建速度: " + time.Now().Sub(t).String())

//slice 创建速度：3.1410928s
//list 创建速度: 41.4315308s
//对于1亿条数据，slice 的创建和添加元素的速度约是list的13倍。
//============================================================
//sli := make([]int, 10)
//for i := 0; i < 1*100000*1000; i++ {
//	sli = append(sli, 1)
//}

//l := list.New()
//for i := 0; i < 1*100000*1000; i++ {
//	l.PushBack(1)
//}
//// 比较遍历
//t := time.Now()
//for _, _ = range sli {
//	//fmt.Printf("values[%d]=%d\n", i, item)
//}
//fmt.Println("遍历slice的速度:" + time.Now().Sub(t).String())
//t = time.Now()
//for e := l.Front(); e != nil; e = e.Next() {
//	//fmt.Println(e.Value)
//}
//fmt.Println("遍历list的速度:" + time.Now().Sub(t).String())

//遍历slice的速度:65.1759ms
//遍历list的速度:28.7595276s
//这差距。。
//============================================================
//    sli:=make([]int ,10)
//    for i := 0; i<1*100000*1000;i++  {
//    sli=append(sli, 1)
//    }

//    l:=list.New()
//    for i := 0; i<1*100000*1000;i++  {
//        l.PushBack(1)
//    }
//    //比较插入
//    t := time.Now()
//    slif:=sli[:100000*500]
//    slib:=sli[100000*500:]
//    slif=append(slif, 10)
//    slif=append(slif, slib...)
//    fmt.Println("slice的插入速度" + time.Now().Sub(t).String())

//    var em *list.Element
//    len:=l.Len()
//    var i int
//    for e := l.Front(); e != nil; e = e.Next() {
//        i++
//        if i ==len/2 {
//            em=e
//            break
//        }
//    }
//    //忽略掉找中间元素的速度。
//    t = time.Now()
//    ef:=l.PushBack(2)
//    l.MoveBefore(ef,em)
//    fmt.Println("list的插入速度: " + time.Now().Sub(t).String())

//slice的插入速度:1.9442905s
//list的插入速度:2.0326ms
//list的优势在快速的插入数据
//======================================================================
package group

import (
	//"encoding/gob"
	//"fmt"
	"io/ioutil"
	"log"

	//"os"
	"encoding/json"
	"sync"

	"github.com/mitchellh/mapstructure"
)

//切片list结构
type Servergroup struct {
	list []interface{}
	lock sync.RWMutex
}

//创建一个空list结构
func NewServergroup() *Servergroup {
	value := &Servergroup{}
	return value
}

//将对象添加到列表末尾
func (servergroup *Servergroup) Add(val interface{}) bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->Add Error:", r)
			ret = false
		}
	}()
	servergroup.list = append(servergroup.list, val)
	return ret
}

//在指定索引处向列表中插入元素,i从0起始
func (servergroup *Servergroup) AddAt(i int32, val interface{}) bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->AddAt Error:", r)
			ret = false
		}
	}()
	servergroup.list = append(servergroup.list, 0)
	copy(servergroup.list[i+1:], servergroup.list[i:])
	servergroup.list[i] = val
	return ret
}

//删除列表指定索引处的元素,i从0起始
func (servergroup *Servergroup) RemoveAt(i int32) bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->RemoveAt Error:", r)
			ret = false
		}
	}()
	servergroup.list = append(servergroup.list[:i], servergroup.list[i+1:]...)
	return ret
}

//删除列表中的最前的一个元素
func (servergroup *Servergroup) RemoveFirst() bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->RemoveFirst Error:", r)
			ret = false
		}
	}()
	servergroup.list = servergroup.list[1:]
	return ret
}

//删除列表中的最后的一个元素
func (servergroup *Servergroup) RemoveLast() bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->RemoveLast Error:", r)
			ret = false
		}
	}()
	servergroup.list = servergroup.list[:len(servergroup.list)-1]
	return ret
}

//删除列表中的所有元素
func (servergroup *Servergroup) Clear() bool {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->Clear Error:", r)
			ret = false
		}
	}()
	servergroup.list = make([]interface{}, 0)
	return ret
}

//按索引获取元素
func (servergroup *Servergroup) GetAt(i int32) (interface{}, bool) {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->GetAt Error:", r)
			ret = false
		}
	}()
	return servergroup.list[i], ret
}

//按范围获取元素
func (servergroup *Servergroup) GetRange(begin int32, end int32) ([]interface{}, bool) {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->GetRange Error:", r)
			ret = false
		}
	}()
	return servergroup.list[begin : end+1], ret
}

//获取列表中的第一个元素
func (servergroup *Servergroup) GetFirst() (interface{}, bool) {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->GetFirst Error:", r)
			ret = false
		}
	}()
	return servergroup.list[0], ret
}

//获取列表中的最后一个元素
func (servergroup *Servergroup) GetLast() (interface{}, bool) {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	ret := true
	defer func() {
		if r := recover(); r != nil {
			log.Println("Servergroup-->GetLast Error:", r)
			ret = false
		}
	}()
	return servergroup.list[len(servergroup.list)-1], ret
}

//统计列表中有多少个数据
func (servergroup *Servergroup) Count() int {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	return len(servergroup.list)
}

//获取列表中的所有值
func (servergroup *Servergroup) Values() []interface{} {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	return servergroup.list
}

//确定元素是否在列表中,只对切片中的值是字符串的有效
func (servergroup *Servergroup) Exists(find interface{}) bool {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	for _, value := range servergroup.list {
		if value == find {
			return true
		}
	}
	return false
}

//返回切片字符串列表,只对切片中的值是字符串的有效
func (servergroup *Servergroup) ToSliceString() []string {
	servergroup.lock.RLock()
	defer servergroup.lock.RUnlock()

	strList := make([]string, len(servergroup.list))
	for k := range servergroup.list {
		val, ok := servergroup.list[k].(string)
		if ok {
			strList[k] = val
		}
	}
	return strList
}

//从文件中读取序列化的数据
func (servergroup *Servergroup) LoadFromFile(filename string) error {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	//	str := string(b)
	//	fmt.Println(str)

	var list []interface{}
	if err2 := json.Unmarshal(b, &list); err2 != nil {
		log.Fatalf("JSON unmarshling failed: %s", err2)
	}

	newlist := NewServergroup()

	for k := range list {

		var servergrouprequest OwlServerGroupRequest
		//将 map 转换为指定的结构体
		if err3 := mapstructure.Decode(list[k], &servergrouprequest); err3 != nil {
			log.Println(err3)
		}

		newlist.Add(servergrouprequest)

	}

	//重新装载数据
	servergroup.list = newlist.list

	return err
}

//将内存数据保存到文件
func (servergroup *Servergroup) SaveToFile(filename string) error {
	servergroup.lock.RLock()

	defer func() {
		servergroup.lock.RUnlock()
		//		if x := recover(); x != nil {
		//			err = fmt.Errorf("Error: Save " + filename)
		//		}
	}()

	data, marshal_err := json.Marshal(servergroup.Values())
	if marshal_err != nil {
		log.Fatalf("Json marshaling failed：%s\n", marshal_err)
	}
	//fmt.Printf("%s\n", data)
	err := ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		log.Fatalf("ioutil WriteFile failed：%s\n", err)
	}

	return err
}
