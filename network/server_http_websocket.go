package network

import (
	//"flag"
	//"io/ioutil"
	//"log"
	"net/http"

	//"os"
	//"fmt"
	"time"

	"github.com/gorilla/websocket"
	// owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
	// owlsystem "github.com/xssed/owlcache/system"
)

//websocket server，此处从server_http.go中轻拆出来方便代码维护

const (
	writeWait  = 7 * time.Second  //发送ping的等待超时时间
	pingPeriod = 2 * time.Second  //ping周期时间，每隔几秒发送一次ping
	pongWait   = 10 * time.Second //等待接收ping返回信息的等待超时时间
)

// 使用websocket配置项
var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//开启websocket
func serveWS(w http.ResponseWriter, r *http.Request) {
	//将Http协议升级为websocket
	client, err := upgrader.Upgrade(w, r, nil)
	//升级失败
	if err != nil {
		owllog.OwlLogWebsocketServer.Info(owltools.JoinString("Upgrade error:", err.Error()))
		return
	}
	//连接成功进来的客户端打印他的IP
	owllog.OwlLogWebsocketServer.Info(owltools.JoinString("WebsocketServer: New client connection:", r.RemoteAddr))

	defer client.Close() //请求结束时资源释放

	//服务端监听pong
	pongHandler(client)

	//定时发送ping命令
	done := make(chan bool)
	defer close(done)
	go pingTicker(client, done, r.RemoteAddr)

	for {

		//监听读取消息
		messageType, payload, err := client.ReadMessage()
		if err != nil {
			owllog.OwlLogWebsocketServer.Info(owltools.JoinString("Read error:", err.Error()))
			break
		}
		//owllog.OwlLogWebsocketServer.Printf("Received message type=%d, payload=\"%s\"\n", messageType, payload)

		//处理接收到的数据
		go WebsocketExe(w, r, string(payload), messageType, client)

	}
}

//服务端监听pong
func pongHandler(client *websocket.Conn) {
	client.SetReadDeadline(time.Now().Add(pongWait)) //设置超时
	client.SetPongHandler(func(string) error {
		//owllog.OwlLogWebsocketServer.Info("Received pong.")
		client.SetReadDeadline(time.Now().Add(pongWait)) //设置超时
		return nil
	})
}

//定时发送ping命令
func pingTicker(client *websocket.Conn, done chan bool, remote_addr string) {
	//定义一个Ticker
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := client.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				owllog.OwlLogWebsocketServer.Info(owltools.JoinString("To ", remote_addr, " sending ping command error:", err.Error()))
			}
		case <-done:
			owllog.OwlLogWebsocketServer.Info(owltools.JoinString("Client:", remote_addr, " stopping ping goroutine."))
			return
		}
	}
}

//Websocket数据执行信息
func WebsocketExe(w http.ResponseWriter, r *http.Request, connstr string, messageType int, client *websocket.Conn) {

	//fmt.Println("WebsocketExe:" + connstr)
	owlhandler := NewOwlHandler()
	owlhandler.owlrequest.WebsocketReceive(w, r, connstr)
	owlhandler.WebsocketHandle(w, r)
	var print []byte
	print = owlhandler.ToWebsocket()

	if err := client.WriteMessage(messageType, print); err != nil {
		owllog.OwlLogWebsocketServer.Info(owltools.JoinString("Write error:", err.Error()))
	}

}
