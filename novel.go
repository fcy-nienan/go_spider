package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
	"os"
	"strings"
)

type Novel struct {
	id       int64
	host     string
	url      string
	name     string
	author   string
	category string
	status   string
	intro    string
	coverUrl string

	chapterAll  ChapterAll
	chapterList []Chapter
	syncDirPath string
}
type Chapter struct {
	url     string
	title   string
	content string
	seq     int
	novelId int
}

type ChapterAll struct {
	url string
}

func (novel *Novel) initDir() {
	err := os.MkdirAll(novel.syncDirPath, 0755) // 权限设置为 rwxr-xr-x
	if err != nil {
		fmt.Println("创建目录失败:", err)
		return
	}
}

func (novel *Novel) parse() {
	client := resty.New()
	resp, err := client.R().Get(novel.url)
	if err != nil {
		log.Printf("NovelParse: 小说：【%s】主页请求失败：%v", novel.name, err)
		return
	}
	html := resp.String()
	host := resp.Request.RawRequest.Host

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("NovelParse: 小说：【%s】主页解析失败：%v", novel.name, err)
		return
	}
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dt").Each(func(i int, s *goquery.Selection) {
		novel.name = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dd:nth-child(2) > span:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		novel.author = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dd:nth-child(2) > span:nth-child(2)").Each(func(i int, s *goquery.Selection) {
		novel.category = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.book_box > dl > dd:nth-child(3) > span:nth-child(1)").Each(func(i int, s *goquery.Selection) {
		novel.status = s.Text()
	})
	doc.Find("body > div.books > div.book_info > div.cover > img").Each(func(i int, s *goquery.Selection) {
		novel.coverUrl = s.Text()
	})
	doc.Find("body > div.books > div.book_about > dl > dd").Each(func(i int, s *goquery.Selection) {
		novel.intro = s.Text()
	})
	novel.syncDirPath = novel.syncDirPath + "\\" + novel.name
	novel.host = host

	var allChapterUrl string
	doc.Find("body > div.books > div.book_more > a").Each(func(i int, s *goquery.Selection) {
		allChapterUrl, _ = s.Attr("href")
	})
	allChapterUrl = "https://" + host + allChapterUrl
	novel.chapterAll = ChapterAll{url: allChapterUrl}
	novel.chapterAll.parse(novel)
}
func (chapterAll *ChapterAll) parse(novel *Novel) (result bool) {
	client := resty.New()
	resp, err := client.R().Get(chapterAll.url)
	if err != nil {
		log.Printf("ChapterAllParse: 小说：【%s】请求章节列表失败：%v", novel.name, err)
		return false
	}
	html := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("ChapterAllParse: 小说：【%s】解析章节列表DOM失败：%v", novel.name, err)
		return false
	}

	doc.Find("body > div.book_last > dl > dd > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		title := s.Text()
		chapterUrl := "https://" + novel.host + href
		novel.chapterList = append(novel.chapterList, Chapter{url: chapterUrl, title: title, seq: i})
	})
	return true
}

func (chapter *Chapter) parse() {
	client := resty.New()
	resp, err := client.R().Get(chapter.url)
	if err != nil {
		log.Printf("ChapterParse: 请求章节【%s】失败：%v", chapter.title, err)
		return
	}
	html := resp.String()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Printf("ChapterParse: 解析章节【%s】DOM失败：%v", chapter.title, err)
		return
	}
	doc.Find("#chaptercontent").Each(func(i int, s *goquery.Selection) {
		chapter.content, _ = s.Html()
	})
}

func (chapter *Chapter) syncFile(novel *Novel) {
	err := os.WriteFile(novel.syncDirPath+"\\"+chapter.title+".txt", []byte(chapter.content), 0644)
	if err != nil {
		log.Printf("ChapterSave: 保存章节【%s】失败：%v", chapter.title, err)
	}
	log.Printf("ChapterSave: 保存章节【%s】成功！", chapter.title)
}
