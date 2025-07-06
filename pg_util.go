package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDatabase() {
	dsn := "user=fcy password=521314 host=127.0.0.1 port=5432 dbname=fcy_db sslmode=disable"
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err.Error())
	}
}
func FindNovelById(id int) (novel Novel) {
	query := "SELECT * from novel where id=$1"
	rows, err := db.Query(query, id)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	novel = Novel{}
	for rows.Next() {
		var id int64
		var name, author, category, status, intro, coverUrl, url string
		err := rows.Scan(&id, &name, &author, &category, &status, &intro, &coverUrl, &url)
		if err != nil {
			panic(err.Error())
		}
		novel.id = id
		novel.name = name
		novel.author = author
		novel.category = category
		novel.status = status
		novel.intro = intro
		novel.coverUrl = coverUrl
		novel.url = url
	}
	return novel
}
func InsertNovel(novel Novel) (id int64) {
	// 插入数据
	insertSQL := "INSERT INTO novel (name, author, category, status, intro, cover_url, url) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	rows, err := db.Query(insertSQL, novel.name, novel.author, novel.category, novel.status, novel.intro, novel.coverUrl, novel.url)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
	}
	return id
}
func InsertChapter(chapter Chapter) (id int64) {
	insertSQL := "insert into chapter (novel_id, seq, title, url, content) values ($1, $2, $3, $4, $5) returning id"
	rows, err := db.Query(insertSQL, chapter.novelId, chapter.seq, chapter.title, chapter.url, chapter.content)
	if err != nil {
		panic(err.Error())
	}
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
	}
	return id
}
