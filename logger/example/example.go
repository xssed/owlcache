package main

import (
	"bufio"
	//"log"
	"os"

	"github.com/xssed/logger"
)

func main() {
	//cutFileHandler := logger.NewCutFileHandler("test", "owl.log", 10*1024*1024) //10M
	cutFileHandler := logger.NewFileHandler("test", "owl.log")

	logger.SetHandlers(logger.Console, cutFileHandler)

	defer logger.Close()

	logger.SetFlags(0) //log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile

	//logger.SetLevel(logger.INFO)

	for i := 0; i < 1000; i++ {
		go func(num int) {
			count := 0
			for {
				logger.Info("Info", num, "-", count)
				logger.Error("Error", num, "-", count)
				logger.Warn("Warn", num, "-", count)
				//logger.Fatal("Fatal", num, "-", count)
				count++
			}
		}(i)
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

}
