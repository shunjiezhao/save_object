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
	hash := meta.Hash
	if hash == "" { //删除标记
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hash = url.PathEscape(hash)             // 转义
	stream, e := GetStream(hash, meta.Size) // 得到从 dataServer请求回来的数据
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, e = io.Copy(w, stream)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	stream.Close()

	w.WriteHeader(http.StatusOK)
}
