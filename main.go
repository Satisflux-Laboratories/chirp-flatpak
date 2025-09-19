package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
)

func main() {

	website := "https://archive.chirpmyradio.com/download?stream=next"

	var wg sync.WaitGroup

	collector := colly.NewCollector(colly.AllowedDomains("archive.chirpmyradio.com"))
	var finalURL string
	var wheel string

	collector.OnRequest(func(r *colly.Request) {
		wg.Add(2)
	})

	collector.OnHTML("a[href]", func(h *colly.HTMLElement) {
		href := h.Attr("href")
		if strings.HasSuffix(href, ".whl") {
			wheel = href
			wg.Done()
		}
	})

	collector.OnResponse(func(r *colly.Response) {
		finalURL = r.Request.URL.String()
		wg.Done()
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Fatalln("Error while scarpig:", err)
	})

	collector.Visit(website)
	wg.Wait()

	if os.Args[1] == "--fullURL" {
		fmt.Print(finalURL)
	} else if os.Args[1] == "--wheelName" {
		fmt.Print(wheel)
	}
}
