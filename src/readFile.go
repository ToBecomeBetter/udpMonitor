package main

import (
	"bufio"
	"fmt"
	"os"
)

// read data from file
func fileModel(fileName, node string, ch chan string) {
	s := "reading file:" + fileName
	fmt.Println(s)
	o, err := os.Open(fileName)
	checkErr(s, err)
	defer o.Close()
	buf := bufio.NewReader(o)
	for {
		l, _ := buf.ReadString('\n')
		if l == "" {
			break
		}
		ch <- l
	}

}
