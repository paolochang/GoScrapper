package main

import (
	"log"
	"net/http"

	"github.com/puerkitobio/goquery"
)

var baseURL string = "https://www.linkedin.com/jobs/search/?keywords=react&start=25"

// $ go get github.com/PuerkitoBio/goquery
func main() {
	var pages = getPages()
	doc, err := goquery.NewDocumentFromReader(res.Body)
}

func getPages() int {
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)
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