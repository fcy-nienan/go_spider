package main

import (
	"fmt"
	"github.com/fcy-nienan/go_mq/mq_client"
	"github.com/fcy-nienan/go_mq/mq_server"
	"strconv"
	"time"
)

func MqTest() {
	go mq_server.StartServer("127.0.0.1:18888")

	time.Sleep(4 * time.Second)

	go func() {
		for {
			for _, value := range mq_server.Qs.Topics {
				fmt.Println("剩余：" + strconv.Itoa(len(value.Messages)) + "条消息！")
			}
			time.Sleep(4 * time.Second)
		}
	}()

	for i := 0; i < 4; i++ {
		go func() {
			producer := mq_client.Client{Address: "127.0.0.1:18888"}
			producer.ConnectServer()
			for i := 0; i < 12; i++ {
				msgStr := "http://www.baidu.com"
				producer.Send("url", []byte(msgStr))
				fmt.Println("生产者发送了一条消息：" + msgStr)
				time.Sleep(100 * time.Millisecond)
			}
		}()
	}

	for i := 0; i < 1; i++ {
		go func() {
			consumer := mq_client.Client{Address: "127.0.0.1:18888"}
			consumer.ConnectServer()
			count := 0
			for {
				msg := consumer.Receive("url")
				fmt.Println(msg)
				if len(msg) == 0 || msg == nil || string(msg) == "EOF" {
					time.Sleep(5 * time.Second)
					continue
				}
				count++
				fmt.Printf("消费者消费了%d条消息:%s\r\n", count, msg)
			}
		}()
	}

	time.Sleep(100000 * time.Second)
}
