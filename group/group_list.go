package group

import (
	"bytes"
	"encoding/json"

	"io/ioutil"
	"log"
	"sync"

	"github.com/xssed/owlcache/tools"
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
func (servergroup *Servergroup) LoadFromFile(folder, filename string) error {
	servergroup.lock.Lock()
	defer servergroup.lock.Unlock()

	b, err := ioutil.ReadFile(folder + filename)
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

	//log.Println(list)

	for k := range list {

		//fmt.Println(tools.Typeof(list[k]))
		temp := list[k].(map[string]interface{})

		var servergrouprequest OwlServerGroupRequest

		servergrouprequest.Cmd = GroupCommandType(temp["Cmd"].(string))
		servergrouprequest.Address = temp["Address"].(string)
		servergrouprequest.Pass = temp["Pass"].(string)
		servergrouprequest.Token = temp["Token"].(string)

		newlist.Add(servergrouprequest)

	}

	//重新装载数据
	servergroup.list = newlist.list

	return err
}

//将内存数据保存到文件
func (servergroup *Servergroup) SaveToFile(folder, filename string) error {
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
	//格式化数据格式
	var out bytes.Buffer
	fmt_err := json.Indent(&out, data, "", "\t")
	if fmt_err != nil {
		log.Fatalf("Json Format failed：%s\n", fmt_err)
	}
	//fmt.Println(out.String())

	//创建文件
	servergroup_dbfile, create_err := tools.CreateFolderAndFile(folder, filename)
	if create_err != nil {
		log.Fatalf("Create File failed：%s\n", create_err)
	}
	_, err := servergroup_dbfile.Write([]byte(out.String()))
	if err != nil {
		log.Fatalf("Write File failed：%s\n", err)
	}
	//释放资源
	servergroup_dbfile.Close()

	//err := ioutil.WriteFile(filename, data, 0777)

	return err
}
