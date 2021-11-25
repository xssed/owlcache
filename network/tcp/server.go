package tcp

import (
	"net"

	owllog "github.com/xssed/owlcache/log"
	owltools "github.com/xssed/owlcache/tools"
)

type server struct {
	address                  string // Address to open connection: ip:port
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

func (s *server) Listen() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		owllog.OwlLogRun.Info("TCPserver: Error starting TCP server.")
	}

	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
		s.onNewClientCallback(client)
	}
}

func New(address string) *server {
	owllog.OwlLogRun.Info(owltools.JoinString("TCPserver: Creating TCP server with address", address))
	server := &server{
		address: address,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

//日志打印
func (s *server) Log(message string) {
	owllog.OwlLogRun.Info(message)
}

//当客户端连接进来时
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

//当客户端链接关闭时
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

//当接收到新消息时
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}
