package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
)

func main() {
	workingDir, outputDir, numWorkers := ParseFlags()

	logFile := setupLogging()
	defer logFile.Close()

	server := StartLocalServer(workingDir)
	defer server.Close()

	VisitAndScrapePage(server.URL, outputDir, numWorkers)
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
	var workers int

	flag.StringVar(&outputDir, "o", "./results", "set the output .csv file directory")
	flag.StringVar(&outputDir, "outputDir", "./results", "set the output .csv file directory")
	flag.StringVar(&workingDir, "d", "./test data", "set the working directory")
	flag.StringVar(&workingDir, "dir", "./test data", "set the working directory")
	flag.IntVar(&workers, "w", 1, "number of goroutines to process files concurrently (100>=w>0)")
	flag.IntVar(&workers, "workers", 1, "number of goroutines to process files concurrently (100>=w>0)")

	flag.Parse()
	if workers <= 0 {
		workers = 1
	} else if workers > 100 {
		workers = 100
	}
	return workingDir, outputDir, workers
}
