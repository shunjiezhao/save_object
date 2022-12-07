package temp

import (
	"net/http"
	"os"
	"strings"
)

// 移除2个暂时文件
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimSuffix(r.URL.EscapedPath(), "\n")
	infoFile := InfoFileName(uuid)
	dataFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(dataFile)
}
