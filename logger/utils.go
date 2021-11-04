package logger

import (
	"fmt"
	"os"
)

//工具函数
//获取文件大小
func fileSize(file string) int64 {
	// 获取文件信息
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

//判断文件是否存在
func isExist(path string) bool {
	// 获取文件信息
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
