package main

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func VisitAndScrapePage(targetURL string, outputDir string, jobs int) {
	log.Println("Visit Web: ", targetURL)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, jobs) // handle the number of go routines

	// scrape all the hyperlinks from the main page, then visit all in parallel
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		subpage := fmt.Sprintf("%s/%s", targetURL, link)

		wg.Add(1)
		ch <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-ch }()
			scrapeSubpage(url, outputDir)
		}(subpage)
	})

	if err := c.Visit(targetURL); err != nil {
		log.Fatal("Visit Web Error: ", err)
	}
	wg.Wait()
	log.Println("Leave Web: ", targetURL)
}

func scrapeSubpage(pageURL string, outputDir string) {
	log.Println("Start scrape: ", pageURL)

	var (
		code []string
		date []string
		body [][]string
	)

	// select specific table and scrape table data
	c := colly.NewCollector()
	const tableSelector = ".matrix.table.table-sm.table-less-padding.table-borderless.table-striped"
	c.OnHTML(tableSelector, func(e *colly.HTMLElement) {
		// scrape head
		e.ForEach("thead tr", func(_ int, h *colly.HTMLElement) {
			h.ForEach("th", func(_ int, th *colly.HTMLElement) {
				txt := strings.Fields(th.Text)
				if len(txt) >= 3 {
					code = append(code, txt[0])
					date = append(date, txt[1]+" "+txt[2])
				}
			})
		})

		// scrape body data
		e.ForEach("tbody tr", func(rows int, b *colly.HTMLElement) {
			if len(body) <= rows {
				body = append(body, []string{})
			}
			b.ForEach("td", func(_ int, td *colly.HTMLElement) {
				attr := td.Attr("class")
				txt := strings.TrimSpace(td.Text)
				newData := mapDataByAttr(attr, txt)
				body[rows] = append(body[rows], newData)
			})
		})
	})

	if err := c.Visit(pageURL); err != nil {
		log.Fatal("Visit Error: ", err)
	}

	// build csv file with scraped data
	BuildCSVFile(pageURL, outputDir, body, code, date)

	log.Println("  end scrape: ", pageURL)
}

func mapDataByAttr(attr string, rawData string) string {
	var newData string

	switch attr {
	// model name
	case "boardmodel text-ellipsis":
		newData = rawData
	// fail
	case "bucket cell-full bg-danger":
		newData = "x"
	// pass
	case "bucket cell-full ":
		if len(rawData) > 0 {
			newData = "1"
		}
	// blank, ignore
	default:
		newData = ""
	}

	return newData
}
