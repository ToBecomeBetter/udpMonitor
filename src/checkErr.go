package main

import "fmt"

func checkErr(info string, err error) {
	if err != nil {
		fmt.Println(info, err)
	}
}
