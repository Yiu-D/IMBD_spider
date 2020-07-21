package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func HttpGetDB(url string) (result string, err error) {
	resp, err1 := http.Get(url)
	if err1 != nil {
		err = err11
		return
	}
	defer resp.Body.Close()
	buf := make([]byte, 4096)
	for {
		n, err2 := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		if err2 != nil && err2 != io.EOF {
			err = err2
			return
		}
		result += string(buf[:n])
	}
	return
}

func SpiderPage(idx int, page chan int) {
	url := "http://www.imdb.cn/imdb250/" + strconv.Itoa((idx))
	//封装HTTPGETDB 爬取 URL对于页面
	result, err := HttpGetDB(url)
	if err != nil {
		fmt.Println("HttpGetDB err:", err)
		return
	}

	//正则取电影名称
	ret := regexp.MustCompile(`<p class="bb">(?s:(.*?))</p>`)
	fileName := ret.FindAllStringSubmatch(result, -1)

	//评分
	ret2 := regexp.MustCompile(`<span><i>(?s:(.*?))</i></span>`)
	fileMark := ret2.FindAllStringSubmatch(result, -1)

	saveFiles(idx, fileName, fileMark)
	page <- idx
}

func saveFiles(idx int, fileName, fileMark [][]string) {
	path := "E:/img/" + "第" + strconv.Itoa(idx) + "页.txt"
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("os.create err:", err)
		return
	}
	defer f.Close()
	n := len(fileName)
	f.WriteString("电影名称" + "\t\t\t" + "评分" + "\n")
	for i := 0; i < n; i++ {
		f.WriteString(fileName[i][1] + "\t\t\t" + fileMark[i][1] + "\n")
	}

}

func toWork(start, end int) {
	fmt.Printf("正在爬取%d 到 %d页...\n", start, end)
	page := make(chan int)
	for i := start; i <= end; i++ {
		go SpiderPage(i, page)
	}
	for i := start; i <= end; i++ {
		fmt.Printf("第 %d 页爬取完毕\n", <-page)
	}

}

func main() {
	var start, end int
	fmt.Print("请输入爬取的起始页(>=1):")
	fmt.Scan(&start)
	fmt.Printf("请输入爬取的结束页(>=%d):", start)
	fmt.Scan(&end)
	toWork(start, end)
}
