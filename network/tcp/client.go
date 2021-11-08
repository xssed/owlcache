package tcp

import (
	"bufio"
	"io"
	"net"

	owllog "github.com/xssed/owlcache/log"
)

type Client struct {
	conn   net.Conn
	Server *server
}

func (c *Client) listen() {

	owllog.OwlLogRun.Info("TCPserver: New connection from ", c.conn.RemoteAddr().String())
	defer c.conn.Close()
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()

			if err == io.EOF {
				owllog.OwlLogRun.Info("TCPserver: Client closed the connection ", c.conn.RemoteAddr().String())
			} else {
				owllog.OwlLogRun.Info("TCPserver: Some problem with reading from client ", c.conn.RemoteAddr().String())
			}
			owllog.OwlLogRun.Info("TCPserver: Done serving ", c.conn.RemoteAddr().String())

			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}

}

func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() error {
	return c.conn.Close()
}
