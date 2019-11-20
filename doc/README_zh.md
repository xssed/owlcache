<a href="https://github.com/xssed/owlcache" target="_blank">English</a> | 中文简介

<div align="center">

# 🦉owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/owl.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)

</div>

 🦉owlcache 是一款由Go编写的轻量级、高性能、无中心分布式的Key/Value内存缓存型的数据共享应用(一定场景下可以作为轻量型数据库来使用)。你可以把每一个节点看做独立的Key/Value数据服务。也可以在默认集群方式下，通过集群查询任意一个节点的数据，得到集群中不同节点之间拥有相同Key的数据列表，最新的数据优先展示。如果开启了同步数据功能，集群中的一个节点Key数据被更改，别的节点中这个Key数据也会更新。   


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
* 🔭**同时支持TCP、HTTP/HTTPS两种方式连接**  
* 🍻**支持Memcache、Redis数据对接**  


## 使用文档
- 📝[中文简体](zh/0.directory.md)
- 📝[English](en/0.directory.md)


## 为什么创建这个项目

* 不认可Memcache、Redis(缓存用途)的集群方式  
* 增加缓存的访问方式  
* 轻松的构建一个高效的数据共享集群  

猫头鹰🦉机灵又可爱。它们脑袋的活动范围为270°🦉。      


![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/group.gif?raw=true)


## 设计  

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/works_zh.png?raw=true)



## 如何编译

编译环境要求
* golang >= 1.9

源码下载
* go命令下载(会自动下载依赖库，如果直接下载源码编译会提示类库缺失)
```shell
go get -u github.com/xssed/owlcache
```

进入owlcache主目录执行编译命令
```shell
go build
```

## 运行
* 注意owlcache.conf文件要跟主程序同目录（下面有介绍动态的设置owlcache.conf文件路径参数）。    
* .conf配置文件必须是统一的UTF-8编码。    
* 请先给在配置文件owlcache.conf中设置<Pass>选项。  

Linux
```shell
./owlcache
```
Windows (DOS下)  
* 注意:Windows生产环境部署中发现，查询请求处理后内存释放会相对慢一些，大概半小时之内才会释放完毕，Linux不存在这个问题，这跟Go的内部机制有关。
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

# 简单使用示例
## 单节点获取Key值
* TCP
命令: `get <key>`
~~~shell
get hello
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
<br>

响应结果例子:
~~~shell
{
    "Cmd": "get",
    "Status": 200,
    "Results": "SUCCESS",
    "Key": "hello",
    "Data": "world",
    "ResponseHost": "127.0.0.1:7721",
    "KeyCreateTime": "2019-04-24T18:05:10.9132377+08:00"
}
~~~

## 集群获取Key值
* 假设现在有三个owlcache服务:127.0.0.1:7721、127.0.0.1:7723、127.0.0.1:7725。每个服务中都有一个Key名称叫hello的数据。


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
{
    "Cmd": "get",
    "Status": 200,
    "Results": "SUCCESS",
    "Key": "hello",
    "Data": [
        {
            "Address": "127.0.0.1:7723",
            "Data": "world7723",
            "KeyCreateTime": "2019-04-10T13:43:01.6576413+08:00",
            "Status": 200
        },
        {
            "Address": "127.0.0.1:7721",
            "Data": "world7721",
            "KeyCreateTime": "2019-04-09T17:50:59.458104+08:00",
            "Status": 200
        },
        {
            "Address": "127.0.0.1:7725",
            "Data": "world7725",
            "KeyCreateTime": "2019-04-08T14:32:20.6934487+08:00",
            "Status": 200
        }
    ],
    "ResponseHost": "127.0.0.1:7721",
    "KeyCreateTime": "0001-01-01T00:00:00Z"
}

~~~

每个节点数据都是独立的，集群中重复Key的数据不会被删除(owlcache默认的集群方式)，查询时会得到一个根据时间排序的数据列表，最新数据优先展示。  



## ......更多请参阅文档的详细说明




## 开发与讨论(不接商业合作)
- 联系我📪:xsser@xsser.cc
- 个人主页🛀:https://www.xsser.cc



