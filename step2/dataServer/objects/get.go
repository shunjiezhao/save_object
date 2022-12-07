package objects

import (
	"ch2/dataServer/locate"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	f := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if f == "" {
		log.Println("没有找到", r.URL.EscapedPath())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	sendFile(w, f)
}
func sendFile(w io.Writer, fileName string) {
	file, _ := os.Open(fileName)
	log.Println(fileName)
	defer file.Close()
	io.Copy(w, file)
}
func getFile(name string) string {
	//TODO:将这块改l
	filePath := path2.Join(dataPath, name)
	log.Println("get file:", filePath)
	files, _ := filepath.Glob(filePath + ".*")
	if len(files) != 1 {
		return ""
	}

	file := files[0]
	log.Println("find file name", file)
	h := sha256.New()
	sendFile(h, file)

	//保障 得到的 hash 一致
	// 防止数据在我们这里受损
	d := url.PathEscape(base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil)))))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove ", file)
		locate.Del(name)
		os.Remove(file)
	}
	return file
}
