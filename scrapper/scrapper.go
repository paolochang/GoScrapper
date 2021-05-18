package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL, filename string

// Scrape Indeed by a keyword
func Scrape(keyword string) {
	
	extStart := time.Now();
	fmt.Println(extStart)
	setBaseURL(keyword)
	setFilename(keyword)
	var jobs []extractedJob
	c := make(chan []extractedJob)
	totalPages := getPages(baseURL, 0)
	fmt.Println("totalPages: ",totalPages)

	for page := 0; page <= totalPages; page++ {
		// extractedJobs := getPage(page)
		go getPage(baseURL, page, c)
		// jobs = append(jobs, extractedJobs...)
	}

	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}
	
	// fmt.Println(jobs)
	writeJobs(jobs)
	extEnd := time.Now()
	// fmt.Println(extEnd)
	timeCons := extEnd.Sub(extStart) / 1000000000
	fmt.Printf("Done, extracted: %d jobs in %02d s\n", len(jobs), timeCons)
}

func setBaseURL(keyword string) {
	baseURL = "https://ca.indeed.com/jobs?q=" + keyword + "&l=Toronto+ON"
}

func getBaseURL() string {
	return baseURL
}

func setFilename(keyword string) {
	filename = keyword + "_" + time.Now().Format("01-02-2006") + ".csv"	
}

func GetFilename() string {
	return filename
}

func getPage(url string, page int, mainC chan<- []extractedJob) {

	var jobs [] extractedJob
	c := make(chan extractedJob)

	pageURL := getBaseURL() + "&start=" + strconv.Itoa(page * 10)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, card *goquery.Selection){
		go extractJob(card, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}

	mainC <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
	id, _ := card.Attr("data-jk")
	title := CleanString(card.Find(".title > a").Text())
	location := CleanString(card.Find(".sjcl > .location").Text())
	salary := CleanString(card.Find(".salaryText").Text())
	summary := CleanString(card.Find(".summary").Text())
	c <- extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
		summary: summary,
	}
}

// CleanString cleans a string
func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func getPages(url string, start int) int {

	if (start != 0) {
		url = getBaseURL() + "&start=" + strconv.Itoa(start)
	}

	pages := 0
	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// fmt.Println(doc)
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

func writeJobs(jobs []extractedJob) {
	c := make(chan []string)
	file, err := os.Create(GetFilename())
	checkErr(err)
	
	w := csv.NewWriter(file)
	defer w.Flush()
	defer file.Close() 

	headers := []string{"Apply", "Title", "Location", "Salary", "Summary"}
	
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		go writeCSV(job, c)
	}

	for i := 0; i < len(jobs); i++ {
		// jobSlice := []string{"https://ca.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		jobSlice := <- c
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func writeCSV(job extractedJob, c chan<- []string) {
	applyLink := "https://ca.indeed.com/viewjob?jk="
	c <- []string{applyLink + job.id, job.title, job.location, job.salary, job.summary}
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