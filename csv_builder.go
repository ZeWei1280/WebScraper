package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	alphabetLength = 26
	rowStart       = 9
	separator      = ";"
)

type CSVBuilder struct {
	fileName string
	filePath string
	header   string
	body     string
	formula  string
	dataCols int
	dataRows int
}

func NewCSVBuilder() *CSVBuilder {
	return &CSVBuilder{}
}

func (b *CSVBuilder) AddFileNameFromURL(pageURL string) *CSVBuilder {
	htmlFile := path.Base(pageURL)
	b.fileName = strings.TrimSuffix(htmlFile, ".html")
	return b
}

func (b *CSVBuilder) AddFilePath(outputDir string) *CSVBuilder {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Fail to create directory %s: %v", outputDir, err)
	}
	b.filePath = filepath.Join(outputDir, b.fileName)
	return b
}

func (b *CSVBuilder) AddHeader(code []string, date []string) *CSVBuilder {
	b.header = fmt.Sprintf("%s%s\n%s%s\n",
		separator, b.separateData(code),
		separator, b.separateData(date))

	return b
}

func (b *CSVBuilder) AddBodyAndSummary(body [][]string) *CSVBuilder {
	b.body = ""
	for _, row := range body {
		b.body += fmt.Sprintf("%s\n", b.separateData(row))
	}

	return b.addFormula(len(body), len(body[0]))
}

func (b *CSVBuilder) BuildCSVFile() {
	f, _ := os.Create(b.filePath + ".csv")
	defer f.Close()
	blankRow := b.separateData(make([]string, b.dataCols))

	fmt.Fprintf(f, "%s%s\n%s%s\n%s",
		b.header,
		blankRow,
		b.formula,
		blankRow,
		b.body,
	)

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
		rowEnd := rowStart + rows - 1
		col := num2CSVColumn(i + 1)
		formulaRange := fmt.Sprintf("%s%d:%s%d", col, rowStart, col, rowEnd)

		pass[i] = fmt.Sprintf("=SUM(%s)", formulaRange)
		fail[i] = fmt.Sprintf("=COUNTIF(%s,\"x\")", formulaRange)
		total[i] = fmt.Sprintf("=COUNTA(%s)", formulaRange)
		passRate[i] = fmt.Sprintf("=IF(%s6=0,\"N/A\",%s5/%s6)", col, col, col)
	}

	b.formula = fmt.Sprintf("%s\n%s\n%s\n%s\n",
		b.separateData(fail),
		b.separateData(pass),
		b.separateData(total),
		b.separateData(passRate))

	return b
}

func (b *CSVBuilder) separateData(data []string) string {
	return strings.Join(data, separator)
}

func num2CSVColumn(num int) string {
	var res string
	for num > 0 {
		mod := (num - 1) % alphabetLength
		res = string(rune('A'+mod)) + res
		num = (num - mod) / alphabetLength
	}
	return res
}
