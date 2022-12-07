package locate

import (
	"ch2/lib/rabbitmq"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var objects = make(map[string]int) // 一次扫描 全部放入
var mutex sync.RWMutex

func Locate(name string) (ok bool) {
	log.Println(name)
	mutex.RLock()
	_, ok = objects[name]
	mutex.RUnlock()
	return
}

func Add(name string) {
	mutex.Lock()
	objects[name] = 1
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
		if Locate(object) {
			println(msg.ReplyTo)
			q.Send(msg.ReplyTo, addr)
		}
	}
}

// ./store/objects_%d/objects
func Collections(path string) {
	files, _ := filepath.Glob(path + "/*")
	for _, file := range files {
		hash := filepath.Base(file)
		objects[hash] = 1
	}
}
