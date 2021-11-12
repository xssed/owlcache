package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"unicode"

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

//高效拼接字符串
func JoinString(args ...string) string {
	var args_buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		args_buffer.WriteString(args[i])
	}
	return args_buffer.String()
}

//对浮点数四舍五入-保留小数点后n位
func RoundedFixed(val float64, n int) float64 {
	change := math.Pow(10, float64(n))
	fv := 0.0000000001 + val //对浮点数产生.xxx999999999 计算不准进行处理
	return math.Floor(fv*change+.5) / change
}

//不区分大小写的对比两个字符串是否相等,高效
func CompareNCS(str1 string, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

//字符串首字母变大写
func Ucfirst(str string) string {
	for i, v := range str {
		return JoinString(string(unicode.ToUpper(v)), str[i+1:])
	}
	return str
}

//字符串首字母变小写
func Lcfirst(str string) string {
	for i, v := range str {
		return JoinString(string(unicode.ToLower(v)), str[i+1:])
	}
	return str
}

//高效将字符串切片拼接字符串
func StringSliceJoinToString(args []string) string {
	var args_buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		args_buffer.WriteString(args[i])
	}
	return args_buffer.String()
}

//高效将错误类型切片拼接字符串
func ErrorSliceJoinToString(args []error) string {
	var args_buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		args_buffer.WriteString(args[i].Error())
	}
	return args_buffer.String()
}

//将[]byte数据进行base64编码转换为加密后字符串
//URL和文件名安全方式是标准方式的变体，其输出用于URL和文件名。因为+和/字符是标准Base64字符对URL和文件名编码不安全，变体即使用-代替+，_（下划线）代替/。
func Base64Encode(input []byte, mode string) string {
	if mode == "url" {
		return base64.URLEncoding.EncodeToString(input)
	}
	return base64.StdEncoding.EncodeToString(input)
}

//将base64加密字符串数据进行base64解码为[]byte,转换错误将返回空[]byte
//URL和文件名安全方式是标准方式的变体，其输出用于URL和文件名。因为+和/字符是标准Base64字符对URL和文件名编码不安全，变体即使用-代替+，_（下划线）代替/。
func Base64Decode(input string, mode string) []byte {
	if mode == "url" {
		uDec, err := base64.URLEncoding.DecodeString(input)
		if err != nil {
			return []byte("")
		}
		return uDec
	}
	uDec, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return []byte("")
	}
	return uDec
}
