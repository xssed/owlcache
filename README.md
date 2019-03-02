# OwlCache

![Image text](https://github.com/xssed/owlcache/blob/master/assets/owl.jpg?raw=true)



OwlCache 是一款由Golang编写的高性能、分布式Key/Value内存缓存系统(一定场景下可以作为轻量型数据库来使用)。
[![IMI License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
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



## 开发计划

Version 0.1 实现单机状态基本功能  
Version 0.2 实现集群数据共享
...... 


## 开发与讨论
- 联系我:xsser@xsser.cc
- 个人主页:https://www.xsser.cc

