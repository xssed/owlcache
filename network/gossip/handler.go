package gossip

import (
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/memberlist"
	owlconfig "github.com/xssed/owlcache/config"
	owllog "github.com/xssed/owlcache/log"
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

	bindport, atio_err := strconv.Atoi(bindPort)
	if atio_err != nil {
		owllog.OwlLogRun.Println("The configuration file <Gossipport> option is not a valid number!")
		os.Exit(0)
	}

	hostname, get_hostname_err := os.Hostname()
	if get_hostname_err != nil {
		owllog.OwlLogRun.Println("When starting the gossip service, getting the Hostname failed! Please check the execution permission of owlcache!")
		os.Exit(0)
	}

	//检查密码
	//没有设置密码
	if owlconfig.OwlConfigModel.GossipDataSyncAuthKey == "" {
		owllog.OwlLogRun.Println("Please set a password first.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	//密码长度过低
	if len(owlconfig.OwlConfigModel.GossipDataSyncAuthKey) <= 10 {
		owllog.OwlLogRun.Println("Password must be greater than 10.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}
	//不能是纯数字
	if _, err := strconv.Atoi(owlconfig.OwlConfigModel.GossipDataSyncAuthKey); err == nil {
		owllog.OwlLogRun.Println("Password cannot be only numbers.Set the <GossipDataSyncAuthKey> option in the configuration file " + owlconfig.OwlConfigModel.Configfile + ".")
		os.Exit(0)
	}

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
	owllog.OwlLogRun.Info("Mark : local member", node.Addr, ":", node.Port)
	return nil
}

func (h *Handler) QueueBroadcast(b []byte) {
	//发送数据到集群
	h.broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte(h.Password), b...),
		notify: nil,
	})

}
