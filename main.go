package main

func main() {

	novel := Novel{}
	novel.url = "https://9dc6fac304032b726558.83646.icu/html/118756/"
	novel.sync_dir_path = "D:\\Code\\novel"
	novel.parse()

	novel.init_config()
	for i := 0; i < 3; i++ {
		chapter := novel.chapter_list[i]
		chapter.parse()
		chapter.sync_file(&novel)
	}

}

