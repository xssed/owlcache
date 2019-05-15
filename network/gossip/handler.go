package gossip

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/memberlist"
	"github.com/xssed/owlcache/tools"
)

type Handler struct {
	broadcasts *memberlist.TransmitLimitedQueue
	nodes      []string
	Password   string
}

func NewHandler() *Handler {
	return &Handler{nodes: []string{}, Password: "owlcache"}
}

func (h *Handler) StartService(str_addresslist []string, passWord string, bindAddress string, bindPort string) error {
	//赋值
	h.nodes = str_addresslist

	if len(passWord) != 0 {
		h.Password = passWord
	}

	bindport, _ := strconv.Atoi(bindPort)

	hostname, _ := os.Hostname()
	c := memberlist.DefaultLocalConfig()
	c.Delegate = &delegate{}
	c.Name = hostname + "-" + tools.GetUUIDString()
	c.BindAddr = bindAddress
	c.BindPort = bindport
	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	if len(h.nodes) > 0 {
		_, err := m.Join(h.nodes)
		if err != nil {
			return err
		}
	}
	h.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 2,
	}
	node := m.LocalNode()
	fmt.Printf("Mark : local member %s:%d\n", node.Addr, node.Port)
	return nil
}

func (h *Handler) QueueBroadcast(b []byte) {
	//发送数据到集群
	h.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte(h.Password), b...),
		notify: nil,
	})

}
