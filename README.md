<div align="center">

# 🦉owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/assets/owl.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)

</div>

 🦉owlcache 是一款由Go编写的轻量级、高性能、无中心分布式的Key/Value内存缓存型的数据共享应用(一定场景下可以作为轻量型数据库来使用)。    


## 亮点与功能简述

* 💡跨平台
* 🚀单机超高性能
* ⛓无中心分布式
* 🌈数据并发安全
* 🔍支持数据过期
* 🖥数据落地存储
* 🎨使用简单，操作命令类似Memcache
* 🔭**同时支持TCP、HTTP两种方式连接**
* ⚔️身份认证
* 📝日志记录


## 设计初衷

我不喜欢造轮子，我最早的想法就是实现一个数据共享应用，它可以非常轻松的构建一个高效的数据共享集群。在集群中的数据，它们可以是共同拥有的，也可以是一个节点拥有其它节点随时来获取。集群里面的所有数据首先要是可“共享”的、可“查阅”的数据。

owl是猫头鹰的意思🦉。机灵又可爱🦉。它们脑袋的活动范围为270°🦉。      


![Image text](https://github.com/xssed/owlcache/blob/master/assets/group.gif?raw=true)

## 使用文档
- 📝http://owl.xsser.cc


## 如何编译

编译环境要求
* golang >= 1.9

源码下载
* go命令下载
```shell
go get github.com/xssed/owlcache
```

编译
```shell
go build
```

## 运行
* 注意owlcache.conf文件要跟主程序同目录（下面有介绍动态的设置owlcache.conf文件路径参数）
```shell
owlcache
```

参数help
* 运行时您可以查看使用帮助 
* 注意运行时的配置参数要优先于*.conf文件里的配置参数

```shell
owlcache -help

Welcome to use owlcache. Version:0.1.3
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
        binding local host ip adress. (default "0.0.0.0")
  -log string
        owlcache log file path.[demo:/var/log/] (default "./log_file/")
  -pass string
        owlcache Http connection password. (default "shi!jie9he?ping6")
```

带配置参数运行的例子
```shell
owlcache -config /var/home/owl.conf -host 127.0.0.1 -log /var/log/ -pass 1245!df2A
```

# 简单使用示例
## 获取Key值
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
    "Data": "world"
}
~~~

......更多请参阅文档的详细说明

## 开发计划

Version 0.1 🚲实现单机状态基本功能  
Version 0.2 🏍实现集群数据共享  
Version 0.3 🚕...... 


## 开发与讨论
- 联系我📪:xsser@xsser.cc
- 个人主页🛀:https://www.xsser.cc

## 开源协议
- [![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)

