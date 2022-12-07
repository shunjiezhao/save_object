package objects

import (
	"ch2/lib/utils"
	"net/http"
	"os"
)

/**
ResponseWriter，Request我就不解释了
*/
var (
	dataPath string
	tmpPath  string
)

func SetAddr(addr string) {
	tmpPath, dataPath = utils.SetAddr(addr)
	os.MkdirAll(tmpPath, 0666)
	os.MkdirAll(dataPath, 0666)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
