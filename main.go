package main

import (
	"WebScraper/utils"
)

func main() {
	workingDir, outputDir, concurrency := utils.ParseFlags()

	logFile := utils.SetupLogging()
	defer logFile.Close()

	server := utils.StartLocalServer(workingDir)
	defer server.Close()

	VisitAndScrapePage(server.URL, outputDir, concurrency)
}
