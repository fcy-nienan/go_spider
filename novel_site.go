package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/fcy-nienan/go_mq/mq_client"
	"github.com/go-resty/resty/v2"
	"strconv"
	"strings"
	"time"
)

// https://m.997b84a6.sbs/map/200.html

func ParseMap() {
	var Producer = mq_client.Client{Address: "127.0.0.1:18888"}
	Producer.ConnectServer()

	client := resty.New()

	for i := 1; i < 2; i++ {
		url := "https://m.997b84a6.sbs/map/" + strconv.Itoa(i) + ".html"
		resp, err := client.R().Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		html := resp.String()
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
		if err != nil {
			fmt.Println(err)
			continue
		}
		host := resp.Request.RawRequest.Host

		doc.Find("body > div.wrap.rank > div.block > ul > li > a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			novelUrl := "https://" + host + href
			//fmt.Println("小说主页：" + novelUrl)
			Producer.Send("novel_main_url", []byte(novelUrl))
		})
	}
	fmt.Println("所有小说解析完毕！")
	fmt.Println("生产者任务完成！")
	//for i := 0; i < 4; i++ {
	//	fmt.Println("所有小说解析完毕！")
	//	time.Sleep(3 * time.Second)
	//}
	//runtime.KeepAlive(Producer)
}

//所有小说页面 url ->novel_main_Url
//单个小说页面 <- url  parse url  chapter_url -> chapter_main_url
//所有章节页面
//单个章节页面 <- chapter_url parse
//一个协程解析所有小说页面，将每个小说的url推送到消息队列
//多个协程从消息队列获取小说url，解析、入库、解析所有章节页面，将小说ID和章节url推送都消息队列
//多个协程从消息队列获取章节url，解析、存入文件、入库

func HandlerNovelUrl() {
	var Consumer = mq_client.Client{Address: "127.0.0.1:18888"}
	Consumer.ConnectServer()
	var Producer = mq_client.Client{Address: "127.0.0.1:18888"}
	Producer.ConnectServer()
	for {
		msg := Consumer.Receive("novel_main_url")
		if msg == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		novel := Novel{url: string(msg), syncDirPath: "D:\\Code\\novel"}
		novel.parse()
		//novel.initDir()
		id := InsertNovel(novel)
		fmt.Println("插入小说：" + novel.name)
		//
		for i := 0; i < len(novel.chapterList); i++ {
			chapter := novel.chapterList[i]
			urlMsg := chapter.url + "#####" + strconv.Itoa(int(id))
			Producer.Send("chapter_main_url", []byte(urlMsg))
		}
	}
}
func HandlerChapterUrl() {
	var Consumer = mq_client.Client{Address: "127.0.0.1:18888"}
	Consumer.ConnectServer()
	for {
		msg := Consumer.Receive("chapter_main_url")
		if msg == nil {
			time.Sleep(2 * time.Second)
			continue
		}
		msgStr := string(msg)
		result := strings.Split(msgStr, "###")
		chapterUrl := result[0]
		novelId, _ := strconv.Atoi(result[1])

		chapter := Chapter{url: chapterUrl, novelId: novelId}
		chapter.parse()
		InsertChapter(chapter)
		fmt.Println("插入章节：" + chapter.title)
	}
}

func ErrorProcess() {
	var Producer = mq_client.Client{Address: "127.0.0.1:18888"}
	Producer.ConnectServer()
}
