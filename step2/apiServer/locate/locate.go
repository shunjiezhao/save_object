package locate

import (
	"ch2/lib/rabbitmq"
	"ch2/lib/rs"
	"ch2/lib/types"
	"encoding/json"
	"os"
	"time"
)

// 返回 id -> 地址的映射
func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()
	}()

	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr
	}
	return
}

// 如果 >= DATA_SHARDS 说明可以恢复原来的数据
func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS
}
