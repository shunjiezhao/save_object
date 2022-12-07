package main

import (
	"ch2/dataServer/heartbeat"
	"ch2/dataServer/locate"
	"ch2/dataServer/objects"
	"ch2/dataServer/temp"
	"flag"
	"fmt"
	"log"
	"net/http"
	path2 "path"
	"strings"
)

var (
	port = ""
	Addr = ""
)

func main() {
	flag.StringVar(&port, "port", "9002", "port")
	flag.Parse()
	log.SetFlags(log.Llongfile)
	Addr := "localhost:" + port
	path := "./store" + fmt.Sprintf("/objects_%s/", strings.Split(Addr, ":")[1])
	objects.SetAddr(path)
	temp.SetAddr(path)

	objPath := path2.Join(path, "objects")
	go heartbeat.StartHeartbeat(Addr)
	go locate.StartLocate(Addr)
	go locate.Collections(objPath)
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(Addr, nil))
}
