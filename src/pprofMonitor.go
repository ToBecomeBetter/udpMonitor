package main

import (
	"net/http"
	"runtime"
	"strconv"

	"github.com/golang/glog"
)

// 性能测试
func pprofMonitor() {
	go func() { // check goes
		http.ListenAndServe("192.168.1.140:6060", nil)
	}()
	go func() {
		http.HandleFunc("/goroutines", func(w http.ResponseWriter, r *http.Request) {
			num := strconv.FormatInt(int64(runtime.NumGoroutine()), 10)
			w.Write([]byte(num))
		})
		http.ListenAndServe("192.168.1.140:6060", nil)
		glog.Info("goroutine stats and pprof listen on 6060")
	}()
}
