package main

import (
	"bufio"
	//"log"
	"os"

	"github.com/xssed/owlcache/logger"
)

func main() {
	log := logger.New()

	cutFileHandler := logger.NewCutFileHandler("test", "owl.log", 10*1024*1024) //10M
	//cutFileHandler := logger.NewFileHandler("test", "owl.log")

	log.SetHandlers(logger.Console, cutFileHandler)

	defer log.Close()

	log.SetFlags(0) //log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

	//logger.SetLevel(logger.INFO)

	for i := 0; i < 1000; i++ {
		go func(num int) {
			count := 0
			for {
				log.Info("Info", num, "-", count)
				log.Error("Error", num, "-", count)
				log.Warn("Warn", num, "-", count)
				//logger.Fatal("Fatal", num, "-", count)
				count++
			}
		}(i)
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}
