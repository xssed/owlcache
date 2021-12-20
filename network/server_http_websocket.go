package network

import (
	//"flag"
	//"io/ioutil"
	"log"
	"net/http"

	//"os"
	"time"

	"github.com/gorilla/websocket"
	// owlconfig "github.com/xssed/owlcache/config"
	// owllog "github.com/xssed/owlcache/log"
	// owlsystem "github.com/xssed/owlcache/system"
)

//websocket server，此处从server_http.go中轻拆出来方便代码维护

const (
	writeWait  = 7 * time.Second  //发送ping的等待超时时间
	pingPeriod = 2 * time.Second  //ping周期时间，每隔几秒发送一次ping
	pongWait   = 10 * time.Second //等待接收ping返回信息的等待超时时间
)

var upgrader = websocket.Upgrader{} // 使用默认websocket配置项

//开启websocket
func serveWS(w http.ResponseWriter, r *http.Request) {
	//将Http协议升级为websocket
	client, err := upgrader.Upgrade(w, r, nil)
	//升级失败
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	//连接成功进来的客户端打印他的IP
	log.Println("WebsocketServer: New client connection:" + r.RemoteAddr)

	defer client.Close() //请求结束时资源释放

	//服务端监听pong
	pongHandler(client)

	//定时发送ping命令
	done := make(chan bool)
	defer close(done)
	go pingTicker(client, done, r.RemoteAddr)

	for {
		messageType, payload, err := client.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received message type=%d, payload=\"%s\"\n", messageType, payload)

		if err := client.WriteMessage(messageType, payload); err != nil {
			log.Println("Write error:", err)
			break
		}
	}
}

//服务端监听pong
func pongHandler(client *websocket.Conn) {
	client.SetReadDeadline(time.Now().Add(pongWait)) //设置超时
	client.SetPongHandler(func(string) error {
		log.Println("Received pong.")
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
				log.Println("To "+remote_addr+" sending ping command error:", err)
			}
			log.Println("To " + remote_addr + " sending ping command success.")
		case <-done:
			log.Println("Stopping ping goroutine.")
			return
		}
	}
}
