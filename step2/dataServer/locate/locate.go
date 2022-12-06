package locate

import (
	"ch2/lib/rabbitmq"
	"log"
	"os"
	"strconv"
)

func Locate(name string) bool {
	log.Println(name)
	_, err := os.Stat(name)
	return err == nil || !os.IsNotExist(err)
}

func StartLocate(path, addr string) {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		object, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		log.Println(path)
		if Locate(path + object) {
			println(msg.ReplyTo)
			q.Send(msg.ReplyTo, addr)
		}
	}
}
