package main

import (
	"ch2/dataServer/heartbeat"
	"ch2/dataServer/locate"
	"ch2/dataServer/objects"
	"log"
	"net/http"
)

func main() {
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.HandleFunc("/objects/", objects.Handler)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}
