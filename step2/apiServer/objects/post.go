package objects

import (
	"ch2/apiServer/heartbeat"
	"ch2/apiServer/locate"
	"ch2/lib/es"
	"ch2/lib/rs"
	"ch2/lib/utils"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// 返回一个token 包括 6 个服务地址及其 uuid
func post(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println("parse size :", err.Error())
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("请传入hash")
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	if locate.Exist(url.PathEscape(hash)) {
		err = es.AddVersion(name, hash, size)
		if err != nil {
			log.Println("添加版本错误:", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		return
	}

	ds := heartbeat.ChooseRandomDataServer(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		log.Println("没有足够多的数据服务器")
		w.WriteHeader(http.StatusServiceUnavailable) // 503
		return
	}
	stream, err := rs.NewRSRemsumblePutStream(ds, name, url.PathEscape(hash), size)
	if err != nil {
		log.Println("创建续传流失败:", err.Error())
		w.WriteHeader(http.StatusServiceUnavailable) // 503
		return
	}
	w.Header().Set("Location", "/temp/"+url.PathEscape(stream.ToToken()))
	w.WriteHeader(http.StatusCreated)
}
