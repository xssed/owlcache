Example
```shell
package main

import (
	"fmt"
	"time"

	"github.com/xssed/owlcache/tools/timeout"
)

func main() {
	to := timeout.New()
	to.SetTimeout("lol", time.Second*7)

	for to.CheckTimeout("lol") {
		fmt.Println(to.CheckTimeout("lol"))
		time.Sleep(time.Second * 1)
	}

}
```