package gossip

import (
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
)

type Handler struct {
	broadcasts *memberlist.TransmitLimitedQueue
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) StartService() error {
	hostname, _ := os.Hostname()
	c := memberlist.DefaultLocalConfig()
	c.Delegate = &delegate{}
	c.BindPort = 0
	//c.BindPort = 7723
	c.Name = hostname + "-" + uuid.NewUUID().String()
	m, err := memberlist.Create(c)
	if err != nil {
		return err
	}
	if len(members) > 0 {
		parts := strings.Split(members, ",")
		_, err := m.Join(parts)
		if err != nil {
			return err
		}
	}
	h.broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return m.NumMembers()
		},
		RetransmitMult: 3,
	}
	node := m.LocalNode()
	fmt.Printf("Local member %s:%d\n", node.Addr, node.Port)
	return nil
}

func (h *Handler) QueueBroadcast(b []byte) {
	//发送数据到集群
	h.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})

}
