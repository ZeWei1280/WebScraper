package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func main() {
	workingDir, outputDir, jobs := ParseFlags()

	logFile := setupLogging()
	defer logFile.Close()

	server := StartLocalServer(workingDir)
	defer server.Close()

	VisitAndScrapePage(server.URL, outputDir, jobs)
}

func StartLocalServer(workingDir string) *httptest.Server {
	fileServerHandler := http.FileServer(http.Dir(workingDir))
	log.Println("Set file server path:", workingDir)

	server := httptest.NewServer(fileServerHandler)
	log.Println("Server started at:", server.URL)
	return server
}

func setupLogging() *os.File {
	logFile, err := os.OpenFile("mylog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Fail to open log file: ", err)
	}
	log.SetOutput(logFile)
	return logFile
}

func ParseFlags() (string, string, int) {
	var workingDir string
	var outputDir string
	var jobs int

	flag.StringVar(&outputDir, "o", "./results", "set the output .csv file directory")
	flag.StringVar(&outputDir, "outputDir", "./results", "set the output .csv file directory")
	flag.StringVar(&workingDir, "d", "./test data", "set the working directory")
	flag.StringVar(&workingDir, "dir", "./test data", "set the working directory")
	flag.IntVar(&jobs, "w", 1, "number of goroutines to process files concurrently (100>=w>0)")
	flag.IntVar(&jobs, "jobs", 1, "number of goroutines to process files concurrently (100>=w>0)")

	flag.Parse()
	if jobs <= 0 {
		// If the jobs parameter is less than or equal to 0, set it to 1 to ensure that at least one goroutine is processing
		jobs = 1
	} else if jobs > 100 {
		// If the jobs parameter is greater than 100, limit it to the maximum value of 100 to prevent excessive concurrency
		jobs = 100
	}
	return workingDir, outputDir, jobs
}
