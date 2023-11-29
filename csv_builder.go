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

func BuildCSVFile(pageURL string, outputDir string, body [][]string, code []string, date []string) {
	csvBuilder := &CSVBuilder{}

	csvBuilder.setFileNameFromURL(pageURL)
	csvBuilder.setFilePath(outputDir)
	csvBuilder.setHeader(code, date)
	csvBuilder.setBody(body)
	csvBuilder.setFormula(len(body), len(body[0]))

	csvBuilder.build()
}

func (b *CSVBuilder) setFileNameFromURL(pageURL string) {
	htmlFile := path.Base(pageURL)
	b.fileName = strings.TrimSuffix(htmlFile, ".html")
}

func (b *CSVBuilder) setFilePath(outputDir string) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Fail to create directory %s: %v", outputDir, err)
	}
	b.filePath = filepath.Join(outputDir, b.fileName)
}

func (b *CSVBuilder) setHeader(code []string, date []string) {
	b.header =
		separator + b.separateData(code) +
			separator + b.separateData(date)
}

func (b *CSVBuilder) setBody(body [][]string) {
	b.body = ""
	for _, row := range body {
		b.body += b.separateData(row)
	}
}

func (b *CSVBuilder) build() {
	f, _ := os.Create(b.filePath + ".csv")
	defer f.Close()
	blankRow := b.separateData(make([]string, b.dataCols))

	fmt.Fprintf(f,
		b.header+
			blankRow+
			b.formula+
			blankRow+
			b.body,
	)

}

func (b *CSVBuilder) setFormula(rows int, cols int) *CSVBuilder {
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
		fail[i] = fmt.Sprintf(`=COUNTIF(%s,"x")`, formulaRange)
		total[i] = fmt.Sprintf("=COUNTA(%s)", formulaRange)
		passRate[i] = fmt.Sprintf(`=IF(%s6=0,"N/A",%s5/%s6)`, col, col, col)
	}

	b.formula =
		b.separateData(fail) +
			b.separateData(pass) +
			b.separateData(total) +
			b.separateData(passRate)

	return b
}

func (b *CSVBuilder) separateData(data []string) string {
	return fmt.Sprintf("%s\n", strings.Join(data, separator))
}

type tableCell struct {
	row    int
	col    int
	cellId string
}

//func NewTableCell(row int, col int) *tableCell {
//	return &tableCell{
//		row: row,
//		col: col,
//		cellId: num2CSVColumn(col),
//	}
//}

func num2CSVColumn(num int) string {
	var res string
	for num > 0 {
		mod := (num - 1) % alphabetLength
		res = string(rune('A'+mod)) + res
		num = (num - mod) / alphabetLength
	}
	return res
}
