package temp

import (
	"ch2/dataServer/locate"
	"ch2/lib/utils"
	"encoding/json"
	"github.com/google/uuid"
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
	w.WriteHeader(http.StatusMethodNotAllowed)
}

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

// 移除2个暂时文件
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.TrimSuffix(r.URL.EscapedPath(), "\n")
	infoFile := InfoFileName(uuid)
	dataFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(dataFile)
}

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
func DataFileName(name string) string {
	str := path2.Join(dataPath, name)
	log.Println((str))
	return str
}
