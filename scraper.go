package main

import (
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func VisitAndScrapePage(targetUrl string, outputDir string, concurrency int) {
	log.Println("Visit Web: ", targetUrl)

	wg := sync.WaitGroup{}
	ch := make(chan struct{}, concurrency)

	collector := colly.NewCollector()
	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		subpage := targetUrl + "/" + link

		wg.Add(1)
		ch <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-ch }()
			scrapeSubpage(url, outputDir)
		}(subpage)
	})

	err := collector.Visit(targetUrl)
	if err != nil {
		log.Fatalln("Visit Web Error: ", err)
	}
	wg.Wait()
	log.Println("Leave Web: ", targetUrl)
}

func scrapeSubpage(pageURL string, outputDir string) {
	log.Println("Start scrape: ", pageURL)

	var (
		code []string
		date []string
		body [][]string
	)

	c := colly.NewCollector()
	const tableSelector = ".matrix.table.table-sm.table-less-padding.table-borderless.table-striped"

	c.OnHTML(tableSelector, func(e *colly.HTMLElement) {
		e.ForEach("thead tr", func(_ int, h *colly.HTMLElement) {
			h.ForEach("th", func(_ int, th *colly.HTMLElement) {
				txt := strings.Fields(th.Text)
				extractCodeAndDate(txt, &code, &date)
			})
		})

		e.ForEach("tbody tr", func(rows int, b *colly.HTMLElement) {
			if len(body) <= rows {
				body = append(body, []string{})
			}
			b.ForEach("td", func(_ int, td *colly.HTMLElement) {
				attr := td.Attr("class")
				txt := strings.TrimSpace(td.Text)
				extractBodyRow(attr, txt, &body[rows])
			})
		})
	})

	err := c.Visit(pageURL)
	if err != nil {
		log.Fatalln("Visit Error: ", err)
	}

	csvBuilder := NewCSVBuilder()
	csvBuilder.
		AddFileNameFromURL(pageURL).
		AddFilePath(outputDir).
		AddBodyAndSummary(body).
		AddHeader(code, date).
		BuildCSVFile()

	log.Println("  end scrape: ", pageURL)
}

func extractBodyRow(attr string, txt string, bodyRow *[]string) {
	newData := mapDataByAttr(attr, txt)
	*bodyRow = append(*bodyRow, newData)
}

func extractCodeAndDate(txt []string, code *[]string, date *[]string) {
	if len(txt) > 0 {
		*code = append(*code, txt[0])
		*date = append(*date, txt[1]+" "+txt[2])
	}
}

func mapDataByAttr(attr string, rawData string) string {
	newData := ""
	switch attr {
	case "boardmodel text-ellipsis":
		newData = rawData
		break
	case "bucket cell-full bg-danger":
		newData = "x"
		break
	case "bucket cell-full ":
		if len(rawData) > 0 {
			newData = "1"
		}
		break
	}
	return newData
}
