package network

import (
	owlconfig "github.com/xssed/owlcache/config"
	"github.com/xssed/owlcache/network/tcp"
)

func startTCP() {

	addr := owlconfig.OwlConfigModel.Host + ":" + owlconfig.OwlConfigModel.Tcpport
	server := tcp.New(addr)

	server.OnNewClient(func(c *tcp.Client) {
		c.Send("Owlcache TCP Link Success...\n")
	})
	server.OnNewMessage(func(c *tcp.Client, message string) {

		//server.Log("OnNewMessage:" + message)//接收到的TCP消息
		owlhandler := NewOwlHandler()
		owlhandler.owlrequest.TCPReceive(message) //解析数据
		owlhandler.TCPHandle()                    //执行数据
		resstr := owlhandler.owlresponse.ConvertToString("TCP") + "\n"
		c.Send(resstr)

	})
	server.OnClientConnectionClosed(func(c *tcp.Client, err error) {

	})

	server.Listen()
}
