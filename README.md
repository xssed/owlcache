English | <a href="https://github.com/xssed/owlcache/blob/master/doc/README_zh.md" target="_blank">中文简介</a>

<div align="center">

# 🦉owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/owl.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)

</div>

 🦉owlcache is a lightweight, high-performance, non-centralized, distributed Key/Value memory-cached data sharing application written by Go (used as a lightweight database in certain scenarios).You can think of each node as a separate Key/Value data service.you can query the data of any node through the cluster to obtain a data list with the same key between different nodes in the cluster. The latest data is displayed first.If the synchronous data function is enabled, the Key data of one node in the cluster is changed, and the Key data of other nodes is also updated.Even if the network environment is not good, many nodes do not synchronize data, you can still query the latest data in the nodes that can be connected.     


## Highlights and features

* 💡Cross-platform operation
* 🚀Single node ultra high performance
* ⛓Non-centralized, distributed
* 🌈Data concurrency security
* 🔍Support data expiration
* 🖥Key/Value Data storage
* 🎨Easy to use, only a few operating commands
* ⚔️Authentication
* 📝Logging
* 🔭**Support both TCP and HTTP/HTTPS connections**  
* 🍻**Support Memcache, Redis data import(String)**  


## Documentation  

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/group.gif?raw=true)

- 📝[中文简体](doc/zh/0.directory.md)
- 📝[English](doc/en/0.directory.md)

</br>

## Design  
 
![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/works_en.png?raw=true)


## How to compile

Compilation environment requirements
* golang >= 1.9

Source download  
* Go command download (will automatically download the dependent library, if you directly download the source code will prompt the class library is missing)  

```shell
go get -u github.com/xssed/owlcache
```

#### ⚠⚠⚠If 'go mod' is ON in your go locale, you need to create a directory locally on your computer, enter the directory, and execute ` git clone https://github.com/xssed/owlcache.git `Command to download the source code.

## Build
* Enter the owlcache home directory and execute the compilation command (in gopath mode, enter the owlcache home directory of gopath directory, and in gomod mode, enter the local directory you created in the previous prompt) 

```shell
go build
```

## Run
* Note that the owlcache.conf file should be in the same directory as the main program.     
* The .conf configuration file must be a uniform UTF-8 encoding.  
* Set the <Pass> option in the configuration file owlcache.conf.     

Linux
```shell
./owlcache
```
Windows (DOS)  
* Note: In the Windows production environment deployment, it is found that the memory release will be relatively slow after the query request processing, and will be released within about half an hour. Linux does not have this problem, which is related to Go's internal mechanism.
```shell
owlcache
```

Parameter help
* You can check out the help before running.
* Note that the runtime configuration parameters take precedence over the configuration parameters in the *.conf file.

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

Example with configuration parameter run
```shell
owlcache -config /var/home/owl.conf -host 127.0.0.1 -log /var/log/ -pass 1245!df2A
```

# Simple use example
## Single node to get the Key value
* TCP
command: `get <key>\n`
~~~shell
get hello\n
~~~

* HTTP
Note: HTTP access data is not verified by password, only other operations that change data require authentication.


|Request parameter        | Parameter value         | 
| ------------- |:-------------: |
| cmd           |  get           | 
| key           |  key name        | 

~~~shell
http://127.0.0.1:7721/data/?cmd=get&key=hello
~~~
<br>

Response result example:
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

## The cluster obtains the Key value.
* Suppose there are now three owlcache services: 127.0.0.1: 7721, 127.0.0.1:7723, 127.0.0.1:7725. Each service has a data called **Key** called **hello**.


|Request parameter        | Parameter value           | 
| ------------- |:-------------: |
| cmd           |  get           | 
| key           |  key name        | 


~~~shell
http://127.0.0.1:7721/group_data/?cmd=get&key=hello
~~~
<br>

Response result example:   
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





## ...more please refer to the detailed description of the document




## Development and discussion(not involved in business cooperation)
- Email📪:xsser@xsser.cc
- Homepage🛀:https://www.xsser.cc



