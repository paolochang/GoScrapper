package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

// var (
// 	limit int = 50
// 	baseURL string = fmt.Sprintf("https://kr.indeed.com/jobs?q=react&limit=%d", limit)
// 	)

// var (
// 	limit int = 25
// 	baseURL string = fmt.Sprintf("https://www.linkedin.com/jobs/search/?keywords=react&start=%d", limit)
// )

// var (
// 	start int = 10
// 	baseURL string = "https://ca.indeed.com/jobs?q=react&l=Toronto+ON&start="
// 	baseURL string = fmt.Sprintf("https://ca.indeed.com/jobs?q=react&l=Toronto+ON&start=%d", start)
// )

var baseURL string = "https://ca.indeed.com/jobs?q=react&l=Toronto+ON"

func main() {
	totalPages := getPages(baseURL, 0)
	fmt.Println(totalPages)

	for page := 0; page <= totalPages; page++ {
		pageURL := baseURL + "&start=" + strconv.Itoa(page * 10)
		fmt.Println("Requesting", pageURL)
	}
}

func getPage(page int) {

}

func getPages(url string, start int) int {

	if (start != 0) {
		url = baseURL + "&start=" + strconv.Itoa(start)
	}

	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	fmt.Println("start: ", start)
	// fmt.Println(doc)
	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection){
		tags := s.Find("a")
		pages = tags.Length()
		// fmt.Println(s.Html())
		tags.Each(func(i int, s *goquery.Selection){
			fmt.Println(s.Attr("aria-label"))
			isNext, _ := s.Attr("aria-label")
			if (isNext == "Next") {
				fmt.Println("YES", isNext, i)
				// getPages(url, i*10)
			}
		})
	})

	return pages
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