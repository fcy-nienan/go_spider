package main

import (
	"fmt"
	"github.com/fcy-nienan/go_mq/mq_server"
	"strconv"
	"time"
)

func main() {
	//MqTest()
	BqgSpider()
}
func BqgSpider() {
	connectDatabase()
	mq_server.StartServer("127.0.0.1:18888")
	time.Sleep(4 * time.Second)

	go func() {
		for {
			for _, value := range mq_server.Qs.Topics {
				fmt.Println("剩余：" + strconv.Itoa(len(value.Messages)) + "条消息！")
			}
			time.Sleep(4 * time.Second)
		}
	}()
	go func() {
		ParseMap()
	}()
	for i := 0; i < 10; i++ {
		go func() {
			HandlerNovelUrl()
		}()
	}

	time.Sleep(100000 * time.Second)
}
