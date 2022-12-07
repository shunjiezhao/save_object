package temp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
	// 1.拿出uuid
	// 2.读取原信息
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := InfoFileName(uuid)
	t, err := readFromFile(infoFile)
	if err != nil {
		log.Println("获取辕信息失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dateFile := t.dataFile()
	file, err := os.OpenFile(dateFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Println("打开文件失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println(uuid, ": 写入失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := file.Stat()
	file.Close()
	if err != nil {
		log.Println(uuid, ": 获取文件元数据失败", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	act := info.Size()
	if act > t.Size {
		os.Remove(dateFile)
		os.Remove(infoFile)
		log.Println("actual size ", act, "exceeds", t.Size)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readFromFile(infoFile string) (*tempInfo, error) {
	b, err := ioutil.ReadFile(infoFile)
	if err != nil {
		return nil, err
	}
	var tmp tempInfo
	json.Unmarshal(b, &tmp)
	return &tmp, nil
}
