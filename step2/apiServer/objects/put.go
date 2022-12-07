package objects

import (
	"ch2/lib/es"
	"ch2/lib/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// 删除的时候 ， hash = “” version + 1
func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("没有设置hash值"))
		return
	}
	size := utils.GetSizeFromHeader(r.Header)
	c, e := storeObject(r.Body, url.PathEscape(hash), size)
	if e != nil {
		log.Println("存储失败", e.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("存储失败"))
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println("版本增加失败: ", e.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(c)
	w.Write([]byte("OK"))
}
