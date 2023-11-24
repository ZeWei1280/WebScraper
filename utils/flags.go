package utils

import "flag"

func ParseFlags() (string, string, int) {
	var workingDir string
	var outputDir string
	var concurrency int

	flag.StringVar(&outputDir, "o", "./results", "set the output .csv file directory")
	flag.StringVar(&outputDir, "outputDir", "./results", "set the output .csv file directory")
	flag.StringVar(&workingDir, "d", "./test data", "set the working directory")
	flag.StringVar(&workingDir, "dir", "./test data", "set the working directory")
	flag.IntVar(&concurrency, "c", 1, "number of goroutines to process files concurrently (100>=c>0)")
	flag.IntVar(&concurrency, "concurrency", 1, "number of goroutines to process files concurrently (100>=c>0)")

	flag.Parse()
	if concurrency <= 0 {
		concurrency = 1
	} else if concurrency > 100 {
		concurrency = 100
	}
	return workingDir, outputDir, concurrency
}
