package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// 移除2个暂时文件
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := InfoFileName(uuid)
	dataFile := infoFile + ".dat"
	err := os.Remove(infoFile)
	if err != nil {
		log.Println("移除文件失败 : %s", err.Error())
	}
	err = os.Remove(dataFile)
	if err != nil {
		log.Println("移除文件失败 : %s", err.Error())
	}
}
