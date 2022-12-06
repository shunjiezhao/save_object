package main

import (
	"ch2/dataServer/heartbeat"
	"ch2/dataServer/locate"
	"ch2/dataServer/objects"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	port = ""
	Addr = ""
)

func main() {
	flag.StringVar(&port, "port", "localhost:9001", "port")
	flag.Parse()
	log.SetFlags(log.Llongfile)
	Addr := "localhost:" + port
	path := "./store" + fmt.Sprintf("/objects_%s/", strings.Split(Addr, ":")[1])
	os.MkdirAll(path, 0666)
	objects.SetAddr(path)
	go heartbeat.StartHeartbeat(Addr)
	go locate.StartLocate(path, Addr)
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe(Addr, nil))
}
