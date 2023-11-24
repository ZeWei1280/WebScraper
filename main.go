package main

import (
	"WebScraper/utils"
)

func main() {
	// flags
	workingDir, outputDir, concurrency := utils.ParseFlags()

	// log
	logFile := utils.SetupLogging()
	defer logFile.Close()

	// local server
	server := utils.StartLocalServer(workingDir)
	defer server.Close()

	// scrape data
	VisitAndScrapePage(server.URL, outputDir, concurrency)
}
