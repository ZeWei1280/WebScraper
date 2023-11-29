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

func (b *CSVBuilder) setFormula(totalRows int, totalCols int) {
	rowEnd := rowStart + totalRows - 1

	// initialize the formula table
	var formula [][]string
	for i := 0; i < 4; i++ {
		formula = append(formula, make([]string, totalCols))
	}
	formula[0][0] = "fail"
	formula[1][0] = "pass"
	formula[2][0] = "total runs"
	formula[3][0] = "pass rate"

	// fill the data to formula table
	for row := 0; row < 4; row++ {
		for col := 0; col < totalCols; col++ {
			tc := NewTableCell(row, col+1)
			if tc.col == 1 {
				// skip the formula name
				continue
			}
			formulaRange := fmt.Sprintf("%s%d:%s%d", tc.cellId, rowStart, tc.cellId, rowEnd)
			formula[0][col] = fmt.Sprintf(`=COUNTIF(%s,"x")`, formulaRange)
			formula[1][col] = fmt.Sprintf("=SUM(%s)", formulaRange)
			formula[2][col] = fmt.Sprintf("=COUNTA(%s)", formulaRange)
			formula[3][col] = fmt.Sprintf(`=IF(%s6=0,"N/A",%s5/%s6)`, tc.cellId, tc.cellId, tc.cellId)
		}
	}
	for _, row := range formula {
		b.formula += b.separateData(row)
	}
}

func (b *CSVBuilder) separateData(data []string) string {
	return fmt.Sprintf("%s\n", strings.Join(data, separator))
}

type TableCell struct {
	row    int
	col    int
	cellId string
}

func NewTableCell(row int, col int) *TableCell {
	return &TableCell{
		row:    row,
		col:    col,
		cellId: num2CSVColumn(col),
	}
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
