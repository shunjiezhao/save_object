package main

import (
	"ch2/apiServer/heartbeat"
	"ch2/apiServer/locate"
	"ch2/apiServer/objects"
	"flag"
	"log"
	"net/http"
	"os"
)

var (
	port = ""
	Addr = ""
)

func main() {
	flag.StringVar(&port, "port", "localhost:9001", "port")
	flag.Parse()
	log.SetFlags(log.Llongfile)
	// api 层 需要了解下面的数据层 有哪些可用，那些不可用
	Addr = "localhost:" + port
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	println(Addr)
	os.MkdirAll("./objects", 0666)

	log.Fatal(http.ListenAndServe(Addr, nil))
}
