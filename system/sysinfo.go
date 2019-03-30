package system

import (
	"fmt"
)

const (
	VERSION      string = "0.1.2"
	VERSION_DATE string = "2018-01-30"
)

func SayHello() {

	fmt.Println("Welcome to use owlcache. Version:" + VERSION + "\nIf you have any questions,Please contact us: xsser@xsser.cc \nProject Home:https://github.com/xssed/owlcache")
	fmt.Println(`                _                _          `)
	fmt.Println(`   _____      _| | ___ __ _  ___| |__   ___ `)
	fmt.Println(`  / _ \ \ /\ / / |/ __/ _' |/ __| '_ \ / _ \`)
	fmt.Println(` | (_) \ V  V /| | (_| (_| | (__| | | |  __/`)
	fmt.Println(`  \___/ \_/\_/ |_|\___\__,_|\___|_| |_|\___|`)
	fmt.Println(`                                            `)

}
