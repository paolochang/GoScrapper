package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var (
	limit int = 50
	baseURL string = fmt.Sprintf("https://kr.indeed.com/jobs?q=react&limit=%d", limit)
	)

// var baseURL string = "https://www.linkedin.com/jobs/search/?keywords=react"

func main() {
	getPages()
}

func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// fmt.Println(doc)
	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection){
		fmt.Println(s.Find("li").Length())
	})

	return 0
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}