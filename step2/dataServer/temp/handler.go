package temp

import (
	"ch2/dataServer/locate"
	"ch2/lib/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	path2 "path"
	"strconv"
	"strings"
)

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

var (
	dataPath string
	tmpPath  string
)

func SetAddr(addr string) {
	tmpPath, dataPath = utils.SetAddr(addr)
}

func (t tempInfo) hash() string {
	s := strings.Split(t.Name, ".")[0]
	log.Println(t.Name, "得到 HASH", s)
	return s
}

func (t tempInfo) id() int {
	s, _ := strconv.Atoi(strings.Split(t.Name, ".")[1])
	log.Println(t.Name, "得到 id", s)
	return s
}

// 将元数据保存文件
func (i *tempInfo) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	b, _ := json.Marshal(i)
	f.Write(b)
	return nil
}

func (i *tempInfo) dataFile() string {
	return path2.Join(tmpPath, i.Uuid+".dat")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPatch {
		patch(w, r)
		return
	}
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodDelete {
		del(w, r)
		return
	}
	if m == http.MethodPost {
		post(w, r)
		return
	}
	if m == http.MethodHead {
		head(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func head(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, err := os.Open(tmpCacheFileName(uuid))
	if err != nil {
		log.Println("打开 缓存中的uuid文件失败：", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Println("获取 缓存中的uuid文件元数据失败：", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
}

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	f, err := os.Open(tmpCacheFileName(uuid))
	if err != nil {
		log.Println("打开 缓存中的uuid文件失败：", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}

func commitTempObject(file string, t *tempInfo) {
	f, _ := os.Open(file)
	d := url.PathEscape(utils.CalculateHash(f))
	f.Close()

	err := os.Rename(file, DataFileName(t.Name)+"."+d+".dat")
	if err != nil {
		log.Println(err)
	}
	locate.Add(t.hash(), t.id())
}

func InfoFileName(uuid string) string {
	str := path2.Join(tmpPath, uuid)
	log.Println(str)
	return str
}
func tmpCacheFileName(uuid string) string {
	str := path2.Join(tmpPath, uuid+".dat")
	log.Println(str)
	return str
}
func DataFileName(name string) string {
	str := path2.Join(dataPath, name)
	log.Println((str))
	return str
}
