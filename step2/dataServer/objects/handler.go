package objects

import (
	"net/http"
)

/**
ResponseWriter，Request我就不解释了
*/
var path string

func SetAddr(addr string) {
	path = addr
}
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r, path)
		return
	}
	if m == http.MethodGet {
		get(w, r, path)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
