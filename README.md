English | <a href="https://github.com/xssed/owlcache/blob/master/doc/README_zh.md" target="_blank">中文简介</a>

<div align="center">

# 🦉owlcache

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/owl.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/owlcache.svg)](https://github.com/xssed/owlcache/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/owlcache.svg?style=popout-square)](https://github.com/xssed/owlcache/releases)

</div>

 🦉owlcache is a lightweight, high-performance, non-centralized, distributed Key/Value in-memory Cache written in Go.It is an independent process and high-performance data middleware, and provides a variety of data get and import methods.You can query a node's key to get all the content with the same key in the node cluster. After operating the Key of a node, the data will be automatically synchronized to all node clusters.       


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
* 🔭**Support TCP and HTTP/HTTPS, WebSocket (Search) connections**  
* 🍻**Support Memcache、 Redis(String)、Url data import**  


## Documentation  

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/group.gif?raw=true)

- 📝[中文简体](doc/zh/0.directory.md)
- 📝[English](doc/en/0.directory.md)


## Sub project     

    🦌 deerfs:Using it, you can build a simple decentralized file system. Project address:<a href="https://github.com/xssed/deerfs" target="_blank"> deerfs</a>     
    
    

## How to compile
<details>
<summary>Show</summary>  

### How to compile
Compilation environment requirements
* golang >= 1.9

Source download  
* Go command download (will automatically download the dependent library, if you directly download the source code will prompt the class library is missing)  

```shell
go get -u github.com/xssed/owlcache
```

#### ⚠⚠⚠If 'go mod' is ON in your go locale, you need to create a directory locally on your computer, enter the directory, and execute ` git clone https://github.com/xssed/owlcache.git `Command to download the source code.

### Build
* Enter the owlcache home directory and execute the compilation command (in gopath mode, enter the owlcache home directory of gopath directory, and in gomod mode, enter the local directory you created in the previous prompt) 

```shell
go build
```

### Run
* Note that the owlcache.conf file should be in the same directory as the main program.     
* The .conf configuration file must be a uniform UTF-8 encoding.  
* Set the <Pass> option in the configuration file owlcache.conf.     

Linux
```shell
./owlcache
```
Windows (DOS)  
* If you plan to use cmd.exe to run owlcache for a long time, please right-click and select [Properties]->[Options]->Close [Quick Edit Mode] and [Insert Mode] in the pop-up menu, otherwise long-running owlcache will appear Caton or dormancy phenomenon.     
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
</details>



## Simple use example
### Single node to get the Key value
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

* Websocket
command: `get <key>`
~~~shell
get hello
~~~

<br>

Response result example:
~~~shell
world
~~~
If it is an HTTP request, there will be Key details in the response message.  
Key: hello  
Keycreatetime: 2021-11-26 18:12:45.1932019 +0800 CST  
Responsehost: 127.0.0.1:7721  


### Single node to get the Key value info
* TCP
command: `get <key> info\n`
~~~shell
get hello info\n
~~~

* HTTP
>owlcache version >= 0.4.2, http no longer supports 'info'

~~`http://127.0.0.1:7721/data/?cmd=get&key=hello&valuedata=info`~~

* Websocket
command: `get <key> info <Custom return string (not necessary, such as UUID)>`
~~~shell
get hello info
~~~
or
~~~shell
get hello info 5c9eff00-3bed-4113-a095-2f3c771683d9
~~~

Response result example:
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

**Attention to the HTTP Status code returned by the HTTP Request. please refer to the "Protocol" chapter.**

### The cluster obtains the Key value.
* Suppose there are now three owlcache HTTP services: 127.0.0.1: 7721, 127.0.0.1:7723, 127.0.0.1:7725. Each service has a data called **Key** called **hello**.


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
world
~~~
The result obtained is the latest value of the update time in the cluster query.


### The cluster obtains the key value information
~~~shell
http://127.0.0.1:7721/group_data/?cmd=get&key=hello&valuedata=info
~~~
<br>

Response result example:   
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
The result is the information about which node in the cluster owns this Key.




## ...more please refer to the detailed description of the document

## Design  
<details>
<summary>Show</summary>   

![Image text](https://github.com/xssed/owlcache/blob/master/doc/assets/works_en.png?raw=true)  

</details>

## Development and discussion(not involved in business cooperation)
- Email📪:xsser@xsser.cc
- Homepage🛀:https://www.xsser.cc



