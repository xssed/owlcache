package gossip

import (
	"github.com/xssed/owlcache/queue"
)

var (
	members = "" //flag.String("members", "", "comma seperated list of members")
	H       = NewHandler()
	Q       = queue.New()
)
