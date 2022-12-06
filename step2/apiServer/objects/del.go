package objects

import (
	"ch2/lib/es"
	"log"
	"net/http"
	"strings"
)

// 拿，找，放
func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	version, err := es.SearchLatestVersion(name)
	if err != nil {
		log.Println("del object : ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = es.PutMetadata(name, version.Version+1, 0, "") // hash 为空表示为删除标记
	if err != nil {
		log.Println("del object : ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
