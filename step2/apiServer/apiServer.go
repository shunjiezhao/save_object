package main

import (
	"ch2/apiServer/heartbeat"
	"ch2/apiServer/locate"
	"ch2/apiServer/objects"
	"log"
	"net/http"
)

func main() {
	// api 层 需要了解下面的数据层 有哪些可用，那些不可用
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
