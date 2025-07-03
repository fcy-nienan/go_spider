package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
)
type Novel struct {
	host string
	url string
	chapter_all ChapterAll
	name      string
	author    string
	category  string
	status    string
	intro     string
	cover_url string
	chapter_list []Chapter

	sync_dir_path string
}
type Chapter struct {
	url string
	title string
	content string
	seq int
}

type ChapterAll struct {
	url string
}
func (novel *Novel) init_config(){
	err := os.MkdirAll(novel.sync_dir_path, 0755) // 权限设置为 rwxr-xr-x
	if err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}
}

func (chapter_all *ChapterAll) parse(novel *Novel) {
	clien := resty.New()
	resp, err := clien.R().Get(chapter_all.url)
	if err != nil {
		fmt.Println(err)
		return;
	}
	html := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))

	doc.Find("body > div.book_last > dl > dd > a").Each(func(i int, s *goquery.Selection) {
		href,_ := s.Attr("href")
		title := s.Text()
		chapter_url := "https://" + novel.host + href
		novel.chapter_list = append(novel.chapter_list, Chapter{ url: chapter_url, title: title, seq: i})
	})
}

func (novel *Novel) parse(){
	clien := resty.New()
	resp, err := clien.R().Get(novel.url)
	if err != nil {
		fmt.Println(err)
		return;
	}
	html := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dt").Each(func(i int, s *goquery.Selection) {
		novel.name = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dd:nth-child(2) > span:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		novel.author = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dd:nth-child(2) > span:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		novel.category = s.Text()
	})

	var all_chapter_url string
	doc.Find("body > div.books > div.book_more > a").Each(func(i int, s *goquery.Selection) {
		all_chapter_url, _ = s.Attr("href")
	})
	host := resp.Request.RawRequest.Host
	all_chapter_url = "https://" + host + all_chapter_url
	novel.host = host
	novel.chapter_all = ChapterAll{url: all_chapter_url}
	novel.sync_dir_path = novel.sync_dir_path + "\\" + novel.name
	novel.chapter_all.parse(novel)
}

func (chapter *Chapter) parse(){
	clien := resty.New()
	resp, err := clien.R().Get(chapter.url)
	if err != nil {
		fmt.Println("请求章节内容失败！【%s】，【%s】", chapter.title, chapter.url)
		return;
	}
	html := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("解析html内容失败！【%s】，【%s】", chapter.title, chapter.url)
		return;
	}
	doc.Find("#chaptercontent").Each(func(i int, s *goquery.Selection) {
		chapter.content, _ = s.Html()
	})
}

func (chapter *Chapter) sync_file(novel *Novel){
	err := os.WriteFile(novel.sync_dir_path + "\\" + chapter.title + ".txt", []byte(chapter.content), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Printf("【%s】章节内容已成功保存\n", chapter.title)
}