package heartbeat

import (
	"ch2/lib/rabbitmq"
	"os"
	"time"
)

func StartHeartbeat(addr string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		q.Publish("apiServers", addr)
		time.Sleep(5 * time.Second)
	}
}
