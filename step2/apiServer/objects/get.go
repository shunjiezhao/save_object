package objects

import (
	"ch2/lib/es"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	var (
		name      = strings.Split(r.URL.EscapedPath(), "/")[2]
		versionId = r.URL.Query()["version"]
		version   = 0
		err       error
	)
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, err := es.GetMetadata(name, version)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" { //删除标记
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	object := url.PathEscape(meta.Hash) // 转义
	stream, e := getStream(object)      // 得到从 dataServer请求回来的数据
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.Copy(w, stream)
}
