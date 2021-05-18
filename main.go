package main

import (
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/paolochang/goscrapper/scrapper"
)

func handleHome(c echo.Context) error {
	return c.File("home.html")
}

func handleScrape(c echo.Context) error {
	defer os.Remove(scrapper.GetFilename())
	keyword := strings.ToLower(scrapper.CleanString(c.FormValue("keyword")))
	scrapper.Scrape(keyword)
	return c.Attachment(scrapper.GetFilename(), scrapper.GetFilename())
}

func main() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrape", handleScrape)
	e.Logger.Fatal(e.Start(":1323"))
}