package locate

import (
	"ch2/lib/rabbitmq"
	"ch2/lib/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int) // 一次扫描 全部放入
var mutex sync.RWMutex

//我们的 Locate 函数不仅要告知某个对象是否存在，同时还需要告知
//本节点保存的是该对象哪个分片
func Locate(name string) int {
	mutex.RLock()
	id, ok := objects[name]
	mutex.RUnlock()
	log.Println(name, id)
	if !ok {
		return -1
	}
	return id
}

func Add(name string, id int) {
	mutex.Lock()
	objects[name] = id
	mutex.Unlock()
}

func Del(name string) {
	mutex.Lock()
	delete(objects, name)
	mutex.Unlock()
}
func StartLocate(addr string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		if id := Locate(object); id != -1 {
			println(msg.ReplyTo)
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: addr, Id: id})
		}
	}
}

// ./store/objects_%d/objects
func Collections(path string) {
	files, _ := filepath.Glob(path + "/*")
	for _, file := range files {
		s := strings.Split(filepath.Base(file), ".") // hash.x.ihash
		i, _ := strconv.ParseInt(s[1], 10, 64)
		log.Println("add file", s[0], " ", i)
		objects[s[0]] = int(i)
	}
}
