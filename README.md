# Owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/assets/owl.jpg?raw=true)



💡 Owlcache 是一款由Golang编写的轻量级、高性能、无中心分布式的Key/Value内存缓存应用(一定场景下可以作为轻量型数据库来使用)。  

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)



## 亮点与功能简述

* 跨平台
* 单机超高性能
* 数据并发安全
* 支持数据过期
* 数据落地存储
* 使用简单，操作命令类似Memcache
* **同时支持TCP、HTTP两种方式连接**
* 身份认证
* 日志记录
* 无中心分布式


## 设计初衷

可以轻松构建一个高效的数据共享与缓存服务集群(偏WEB方向)。Owl是猫头鹰的意思。


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
Welcome to use owlcache. Version:0.1  by:d4rkdu0
Usage of owlcache:
  -config string
        Owlcache config file path.[demo:/var/home/owl.conf] (default "owlcache.conf")
  -host string
        binding local host ip adress. (default "127.0.0.1")
  -log string
        Owlcache log file path.[demo:/var/log/]
  -pass string
        Owlcache connection password. (default "shi!jie9he?ping6")
```

带配置参数运行的例子
```shell
owlcache -config /var/home/owl.conf -host 127.0.0.1 -log /var/log/ -pass 1245!df2A
```


## 使用文档
- http://owl.xsser.cc


## 开发计划

Version 0.1 实现单机状态基本功能  
Version 0.2 实现集群数据共享
...... 


## 开发与讨论
- 联系我:xsser@xsser.cc
- 个人主页:https://www.xsser.cc

## 开源协议
- [![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)

