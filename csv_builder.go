package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type CSVBuilder struct {
	fileName  string
	filePath  string
	separator string
	header    string
	body      string
	summary   string
	dataCols  int
	dataRows  int
}

func NewCSVBuilder() *CSVBuilder {
	return &CSVBuilder{
		separator: ";",
	}
}

func (b *CSVBuilder) AddFileNameFromURL(pageURL string) *CSVBuilder {
	htmlFile := path.Base(pageURL)
	filename := strings.TrimSuffix(htmlFile, ".html")
	b.fileName = filename
	return b
}

func (b *CSVBuilder) AddFilePath(outputDir string) *CSVBuilder {
	os.MkdirAll(outputDir, os.ModePerm)
	b.filePath = filepath.Join(outputDir, b.fileName)
	return b
}

func (b *CSVBuilder) BuildCSVFile() {
	f, _ := os.Create(b.filePath + ".csv")
	defer f.Close()
	blankLine := b.separateData(make([]string, b.dataCols)) + "\n"

	fmt.Fprintf(f, "%s%s%s%s%s",
		b.header,
		blankLine,
		b.summary,
		blankLine,
		b.body)
}

func (b *CSVBuilder) AddHeader(code []string, date []string) *CSVBuilder {
	b.header =
		b.separator + b.separateData(code) + "\n" +
			b.separator + b.separateData(date) + "\n"

	return b
}

func (b *CSVBuilder) AddBodyAndSummary(body [][]string) *CSVBuilder {
	for _, row := range body {
		b.body += b.separateData(row) + "\n"
	}
	return b.addFormula(len(body), len(body[0]))
}

func (b *CSVBuilder) addFormula(rows int, cols int) *CSVBuilder {
	var pass = make([]string, cols)
	var fail = make([]string, cols)
	var total = make([]string, cols)
	var passRate = make([]string, cols)
	pass[0] = "pass"
	fail[0] = "failed"
	total[0] = "total runs"
	passRate[0] = "pass rate"

	for i := 1; i < cols; i++ {
		rowStart := 9
		rowEnd := rowStart + rows - 1
		col := num2CSVColumn(i + 1)
		formulaRange := col + strconv.Itoa(rowStart) + ":" + col + strconv.Itoa(rowEnd)

		pass[i] = "=SUM(" + formulaRange + ")"
		fail[i] = "=COUNTIF(" + formulaRange + ",\"x\")"
		total[i] = "=COUNTA(" + formulaRange + ")"
		passRate[i] = "=IF(" + col + "6=0,\"N/A\"," + col + "5/" + col + "6)"
	}

	b.summary = b.separateData(fail) + "\n" +
		b.separateData(pass) + "\n" +
		b.separateData(total) + "\n" +
		b.separateData(passRate) + "\n"
	return b
}

func (b *CSVBuilder) separateData(data []string) string {
	return strings.Join(data, b.separator)
}

func num2CSVColumn(num int) string {
	var res string
	for num > 0 {
		mod := (num - 1) % 26
		res = string(rune('A'+mod)) + res
		num = (num - mod) / 26
	}
	return res
}
