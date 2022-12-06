package main

import (
	"ch1/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/objects", handler.Handler)
	http.ListenAndServe("localhost:8080", nil)
}
