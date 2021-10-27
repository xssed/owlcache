package config

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/xssed/owlcache/tools"
)

//创建一个全局URL缓存配置变量
var OwlUCConfigModel *Sites

//结构体-站点集合信息
type Sites struct {
	SiteList []Site `xml:"site"`
}

//结构体-站点信息
type Site struct {
	Host            string   `xml:"host"`              //需要映射的URL基础网址。必填项。默认为空。例:127.0.0.1:3306
	Headers         []Header `xml:"header"`            //向映射网址发送GET访问请求时，将会把当前HTTP请求Request headers的子项附带转发到映射网址。默认为空。
	MaxStorageLimit uint64   `xml:"max_storage_limit"` //映射请求返回的数据在内存中允许的最大存储字节数。单位byte。默认最大值得为5242880byte->5M
	CheckToken      int      `xml:"check_token"`       //是否要校验owl的token。默认为0 。0为关闭，1为开启。
	KeyExpire       int      `xml:"key_expire"`        //存储在内存的有效时间。单位秒。默认为0。0为永久。
	Timeout         uint64   `xml:"timeout"`           //映射请求的超时时间。单位毫秒。默认5000ms->5s
	Proxy           string   `xml:"proxy"`             //映射请求的代理服务器。默认为空
}

//结构体-用于指定Request Headers转发
type Header struct {
	Value []string `xml:"value"`
}

//创建一个默认配置文件的实体
func NewDefaultOwlUCConfig() *Sites {
	return &Sites{}
}

//URL缓存初始化加载配置
func OwlUCConfigInit() {

	//创建一个默认URL缓存初始化配置模型
	Temp_OwlUCConfigModel := NewDefaultOwlUCConfig()

	//配置文件路径
	var cpath string = tools.JoinString(OwlConfigModel.DBfile, OwlConfigModel.Urlcache_Filename)

	//打开URL缓存文件载入数据
	ucconfig_content, err := ioutil.ReadFile(cpath)
	//监控错误
	if err != nil {
		fmt.Println(err.Error())
		fmt.Print("Can not read url cache configuration file. now exit\n")
		os.Exit(0)
	}
	//判断文件是否为UTF-8编码
	if !tools.ValidUTF8(cpath) {
		fmt.Print("The url cache configuration file is not UTF-8 encode. now exit\n")
		os.Exit(0)
	}
	//XML配置文件将数据绑定
	xerr := xml.Unmarshal(ucconfig_content, Temp_OwlUCConfigModel)
	if xerr != nil {
		fmt.Println(xerr.Error())
		fmt.Print("url cache XML configuration file data error. now exit\n")
		os.Exit(0)
	}

	//管理配置默认值
	OwlUCConfigModel = OwlUCConfigCheckDefaultValue(Temp_OwlUCConfigModel)

	//fmt.Println(OwlUCConfigModel) //打印出最终赋值后的配置信息

	//执行步骤信息
	fmt.Println("owlcache  url cache configuration initialization is complete...")

}

//管理配置默认值
func OwlUCConfigCheckDefaultValue(sites *Sites) *Sites {

	temp_sites := NewDefaultOwlUCConfig()

	for _, site := range sites.SiteList {
		if site.MaxStorageLimit == 0 {
			site.MaxStorageLimit = 5242880
		}
		if site.Timeout == 0 {
			site.Timeout = 5000
		}
		temp_sites.SiteList = append(temp_sites.SiteList, site)
	}

	return temp_sites

}
