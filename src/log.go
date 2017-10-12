//日志打印

package main

import (
	"log"
	"os"
)

var setPath = logPath()

func sunLog() *log.Logger {
	file, err := os.Create(setPath)
	if err != nil {
		log.Fatalln("fail to create  udp.log file!")
	}
	logger := log.New(file, "[UDP] ", log.Lshortfile|log.LstdFlags)
	return logger
}
