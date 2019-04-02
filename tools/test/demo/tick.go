package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {

	task_serverlistbackup, err := strconv.Atoi("5")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	//	task_serverlistbackup, err := time.ParseDuration("1" + "m")
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		os.Exit(0)
	//	}

	fmt.Println(task_serverlistbackup)

	aa := time.Second * time.Duration(task_serverlistbackup)

	ticker := time.NewTicker(aa)
	go func() {
		for _ = range ticker.C {
			fmt.Printf("ticked at %v", time.Now(), "\n")
		}
	}()

	time.Sleep(time.Minute * 10)
}
