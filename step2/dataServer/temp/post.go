package temp

import (
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 返回 uuid 表示已经做好准备
func post(w http.ResponseWriter, r *http.Request) {
	output := uuid.New().String()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Println("获取size失败")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("内部错误"))
		return
	}
	t := tempInfo{
		Uuid: uuid,
		Name: name,
		Size: size,
	}
	infoFile := InfoFileName(uuid)
	err = t.WriteToFile(infoFile)
	if err != nil {
		log.Printf("uuid:%s name:%s 写入文件失败\n", t.Uuid, t.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("内部错误"))
		return
	}
	f, err := os.Create(infoFile + ".dat")
	if err != nil {
		log.Printf("uuid:%s name:%s 创建文件失败\n", t.Uuid, t.Name)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("内部错误"))
		return
	}
	defer f.Close()
	w.Write([]byte(uuid))
}
