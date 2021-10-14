package tools

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/saintfish/chardet"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//生成UUID
func GetUUIDString() string {
	uuid := uuid.Must(uuid.NewV4())
	return uuid.String()
}

//字符串  客户端IP+PORT转IP
func RemoteAddr2IPAddr(key string) string {
	str := strings.Split(key, ":")
	return str[0]
}

//创建目录
func CreateFolderAndFile(folder, filename string) (*os.File, error) {

	//日志存储目录校验
	temp_last := folder[len(folder)-1:]
	if temp_last != "/" {
		folder = folder + "/"
	}
	//fmt.Println(folder)

	if folder != "" {
		os.MkdirAll(folder, os.ModePerm)
	}
	newfile := folder + filename

	file, err := os.OpenFile(newfile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	return file, err

}

//判断返回传入值的类型
func Typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

//检验文件内容是否为UTF-8
func ValidUTF8(filename string) bool {
	rawBytes, readfile_err := ioutil.ReadFile(filename)
	if readfile_err != nil {
		return false
	}
	detector := chardet.NewTextDetector()
	charset, detectbest_err := detector.DetectBest(rawBytes)
	if detectbest_err != nil {
		return false
	}
	if charset.Charset == "UTF-8" {
		return true
	} else {
		return false
	}
}
