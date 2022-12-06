package handler

import (
	"net/http"
	"path"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		get(w, r)
		return
	}

	if r.Method == "PUT" {
		put(w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}

func put(w http.ResponseWriter, r *http.Request) {
	// url ->/objects/object_name
	println(path.Join("./objects", strings.Split(r.URL.EscapedPath(), "/")[2]))
}

func get(w http.ResponseWriter, r *http.Request) {

}
