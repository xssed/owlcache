<a href="https://github.com/xssed/owlcache" target="_blank">English</a> | 中文简介

<div align="center">

# 🦉owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/owl.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)

</div>

 🦉owlcache 是一款由Go编写的轻量级、高性能、无中心、分布式的Key/Value内存缓存。   


## 亮点与功能简述

* 💡跨平台运行
* 🚀单机超高性能
* ⛓无中心分布式
* 🌈数据并发安全
* 🔍支持数据过期
* 🖥数据落地存储
* 🎨使用简单，操作命令只有几个
* ⚔️身份认证
* 📝日志记录
* 🔭**支持TCP和HTTP/HTTPS、WebSocket(搜索)连接**  
* 🍻**支持Memcache、Redis(String)、Url数据对接**  


## 使用文档

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/group.gif?raw=true)

- 📝[中文简体](zh/0.directory.md)
- 📝[English](en/0.directory.md)

## 如何编译
<details>
<summary>Show</summary>  


编译环境要求
* golang >= 1.9

源码下载
* go命令下载(会自动下载依赖库，如果直接下载源码编译会提示类库缺失)
* go get命令无法执行请检查本机是否安装Git服务和设置Go环境  

```shell
go get -u github.com/xssed/owlcache
```

#### ⚠⚠⚠如果你的Go语言环境开启了`GOMOD`,你需要在电脑本地创建一个目录,进入该目录，再次执行`git clone https://github.com/xssed/owlcache.git`命令将源代码下载. 

### 编译
* 进入owlcache主目录执行编译命令(GOPATH模式下进入GOPATH目录的owlcache主目录，GOMOD模式则进入上一步提示中你自己创建的本地目录)
```shell
go build
```

### 运行
* 注意owlcache.conf文件要跟主程序同目录（下面有介绍动态的设置owlcache.conf文件路径参数）。    
* .conf配置文件必须是统一的UTF-8编码。    
* 请先给在配置文件owlcache.conf中设置<Pass>选项。  

Linux
```shell
./owlcache
```
Windows (DOS)   
* 如果你打算使用cmd.exe长时间运行owlcache，请右键，在弹出菜单中选择【属性】->【选项】->关闭【快速编辑模式】和【插入模式】，否则长时间运行owlcache会出现卡顿或者休眠现象。  
```shell
owlcache
```

参数help
* 运行前您可以查看使用帮助 
* 注意运行时的配置参数要优先于*.conf文件里的配置参数

```shell
owlcache -help
```
```shell
Welcome to use owlcache. Version:XXX
If you have any questions,Please contact us: xsser@xsser.cc
Project Home:https://github.com/xssed/owlcache
                _                _
   _____      _| | ___ __ _  ___| |__   ___
  / _ \ \ /\ / / |/ __/ _' |/ __| '_ \ / _ \
 | (_) \ V  V /| | (_| (_| | (__| | | |  __/
  \___/ \_/\_/ |_|\___\__,_|\___|_| |_|\___|

Usage of owlcache:
  -config string
        owlcache config file path.[demo:/var/home/owl.conf] (default "owlcache.conf")
  -host string
        binding local host ip address. (default "0.0.0.0")
  -log string
        owlcache log file path.[demo:/var/log/] (default "./log_file/")
  -pass string
        owlcache Http connection password. (default "")
```

带配置参数运行的例子
```shell
owlcache -config /var/home/owl.conf -host 127.0.0.1 -log /var/log/ -pass 1245!df2A
```
</details>


## 简单使用示例
### 单节点获取Key值
* TCP
命令: `get <key>\n`
~~~shell
get hello\n
~~~

* HTTP
注意:HTTP获取数据不用密码验证，只有其他更改数据的操作需要验证身份。


|请求参数        | 参数值          | 
| ------------- |:-------------: |
| cmd           |  get           | 
| key           |  key名称        | 

~~~shell
http://127.0.0.1:7721/data/?cmd=get&key=hello
~~~

* Websocket
command: `get <key>`
~~~shell
get hello
~~~

<br>

响应结果例子:
~~~shell
world
~~~
如果是HTTP请求，在响应报文中会有Key的详细信息  
Key: hello  
Keycreatetime: 2021-11-26 18:12:45.1932019 +0800 CST  
Responsehost: 127.0.0.1:7721  


### 单节点获取Key值的信息
* TCP
命令: `get <key> info\n`
~~~shell
get hello info\n
~~~

* HTTP
>owlcache 版本 >= 0.4.2, Http不再支持 'info'

~~`http://127.0.0.1:7721/data/?cmd=get&key=hello&valuedata=info`~~

* Websocket
命令: `get <key> info <Custom return string (not necessary, such as UUID)>`
~~~shell
get hello info
~~~
或者
~~~shell
get hello info 5c9eff00-3bed-4113-a095-2f3c771683d9
~~~

响应结果例子:
~~~shell
{
    "Cmd": "get",
    "Status": 200,
    "Results": "SUCCESS",
    "Key": "hello",
    "Data": "d29ybGQ=",
    "ResponseHost": "127.0.0.1:7721",
    "KeyCreateTime": "2021-11-09T14:12:36.8431596+08:00"
}
~~~

**注意HTTP请求返回的HTTP状态码，解释请参考“通讯协议”章节。**

### 集群获取Key值
* 假设现在有三个owlcache HTTP服务:127.0.0.1:7721、127.0.0.1:7723、127.0.0.1:7725。每个服务中都有一个Key名称叫hello的数据。


|请求参数        | 参数值          | 
| ------------- |:-------------: |
| cmd           |  get           | 
| key           |  key名称        | 


~~~shell
http://127.0.0.1:7721/group_data/?cmd=get&key=hello
~~~
<br>

响应结果例子:   
~~~shell
world
~~~
得到的结果是集群查询中更新时间最新的那一个值。


### 集群获取Key值的信息
~~~shell
http://127.0.0.1:7721/group_data/?cmd=get&key=hello&valuedata=info
~~~
<br>

响应结果例子:   
~~~shell
[
    {
        "Address": "127.0.0.1:7721",
        "Data": "d29ybGQ=",
        "Key": "hello",
        "KeyCreateTime": "2021-11-26T18:12:45.1932019+08:00",
        "Status": 200
    },
    {
        "Address": "127.0.0.1:7723",
        "Data": "d29ybGQ=",
        "Key": "hello",
        "KeyCreateTime": "2021-11-12T11:34:53.0952566+08:00",
        "Status": 200
    },
    {
        "Address": "127.0.0.1:7725",
        "Data": "d29ybGQ=",
        "Key": "hello",
        "KeyCreateTime": "2021-11-11T11:34:53.0952522+08:00",
        "Status": 200
    }
]

~~~
结果是有关集群中哪个节点拥有此Key的信息。
  


## ......更多请参阅文档的详细说明

## 设计  
<details>
<summary>Show</summary> 

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/works_zh.png?raw=true)  

</details>

## 开发与讨论(不接商业合作)
- 联系我📪:xsser@xsser.cc
- 个人主页🛀:https://www.xsser.cc



