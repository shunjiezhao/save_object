package versions

import (
	"ch2/lib/es"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed) // code 405
		return
	}
	from := 0
	size := 1000 // 一次获取多少个
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := es.SearchAllVersions(name, from, size) // 获取所有元信息
		if err != nil {
			log.Println("versions handler is : ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, meta := range metas {
			b, _ := json.Marshal(meta)
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size { // 没有更多了
			return
		}
		from += size // 新的起点
	}

}
