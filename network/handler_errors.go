package network

import (
	"errors"
)

var (
	ErrorKeyNotFound           = errors.New("Not found key in cache")
	ErrorKeyNotExist           = errors.New("Key does not exist")
	ErrorOpenHttpsSelected     = errors.New("Configuration file Open_Https filled in error")
	ErrorCloseTcp              = errors.New("Configuration file CloseTcp filled in error")
	ErrorGroupDataSync         = errors.New("Configuration file GroupDataSync filled in error")
	ErrorCors                  = errors.New("Configuration file Cors filled in error")
	ErrorUCMaxStorageLimitOver = errors.New("The urlcache request succee, but the response data exceed the maximum value of the <max_storage_limit> configuration parameter.")
	ErrorGossipGetUrlData      = errors.New("Gossip get url data error.")
)
