package network

import (
	"errors"
)

var (
	ErrorKeyNotFound       = errors.New("Not found key in cache")
	ErrorKeyNotExist       = errors.New("Key does not exist")
	ErrorOpenHttpsSelected = errors.New("Configuration file Open_Https filled in error")
)
