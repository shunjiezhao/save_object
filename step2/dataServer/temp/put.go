package temp

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// 文件上传成功
func put(w http.ResponseWriter, r *http.Request) {
	// take uuid
	uuid := strings.TrimSuffix(strings.Split(r.URL.EscapedPath(), "/")[2], "\n")
	infoFile := InfoFileName(uuid)
	t, err := readFromFile(infoFile)
	if err != nil {
		log.Println("获取辕信息失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dataFile := t.dataFile()
	f, err := os.Open(dataFile)
	if err != nil {
		log.Println("打开文件失败", dataFile, " ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := f.Stat()
	f.Close()
	if err != nil {
		log.Println(uuid, ": 获取文件元数据失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	act := info.Size()
	os.Remove(infoFile)
	if act != t.Size {
		os.Remove(dataFile)
		log.Println("actual size mismatch, expect", t.Size, "acutal: ", act)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(dataFile, t)
}
