package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

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
	var jobs []extractedJob
	totalPages := getPages(baseURL, 0)
	fmt.Println("totalPages: ",totalPages)

	for page := 0; page <= totalPages; page++ {
		extractedJobs := getPage(page)
		jobs = append(jobs, extractedJobs...)
	}

	// fmt.Println(jobs)
	writeJobs(jobs)
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)
	
	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Apply", "Title", "Location", "Salary", "Summary"}
	
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://ca.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int) []extractedJob {

	var jobs [] extractedJob

	pageURL := baseURL + "&start=" + strconv.Itoa(page * 10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection){
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(card *goquery.Selection) extractedJob {
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title > a").Text())
	location := cleanString(card.Find(".sjcl > .location").Text())
	salary := cleanString(card.Find(".salaryText").Text())
	summary := cleanString(card.Find(".summary").Text())
	return extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
		summary: summary,
	}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
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

	doc.Find(".pagination-list").Each(func(i int, s *goquery.Selection){
		tags := s.Find("a")
		tags.Each(func(i int, s *goquery.Selection){
			// fmt.Println(s.Attr("aria-label"))
			isNext, _ := s.Attr("aria-label")
			if (isNext == "Next") {
				pages = getPages(url, (pages - 1)*10)
			} else {
				if i, err := strconv.Atoi(isNext); err == nil {
					pages = i
				}
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