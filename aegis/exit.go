package aegis

import (
	"fmt"
	"os"
	"os/signal"

	owllog "github.com/xssed/owlcache/log"
)

//捕获程序正常退出操作 ctrl+c
func OnExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	owllog.OwlLogRun.Info("owlcache is stoped") //日志记录
	fmt.Println("owlcache  is stoped \nBye!")
}
