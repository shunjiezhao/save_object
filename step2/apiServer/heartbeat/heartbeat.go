package heartbeat

import (
	"ch2/lib/rabbitmq"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time) // 记录发来心跳消息的时间
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body)) // 双引号去掉
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now() //记录时间
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		select {
		case <-time.After(5 * time.Second):
			mutex.Lock()
			for s, t := range dataServers {
				// 10s已经过去了 超时
				if t.Add(10 * time.Second).Before(time.Now()) {
					delete(dataServers, s)
				}
			}
			mutex.Unlock()
		}
	}
}

func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
