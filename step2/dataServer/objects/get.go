package objects

import (
	"ch2/dataServer/locate"
	"ch2/lib/utils"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	f := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if f == "" {
		log.Println("没有找到", r.URL.EscapedPath())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(f, w)

}
func sendFile(fileName string, w http.ResponseWriter) {
	file, _ := os.Open(fileName)
	log.Println(fileName)
	defer file.Close()
	io.Copy(w, file)
}
func getFile(hash string) string {
	//TODO:将这块改l
	file := path2.Join(dataPath, hash)
	log.Println("get file:", file)

	f, _ := os.Open(file)
	defer f.Close()
	//保障 得到的 hash 一致
	// 防止数据在我们这里受损
	d := url.PathEscape(utils.CalculateHash(f))
	if d != hash {
		log.Println("object hash mismatch, remove ", file)
		locate.Del(hash)
		os.Remove(file)
	}
	return file
}
