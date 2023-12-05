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
	separator      = ";"
)

type CSVBuilder struct {
	fileName string
	filePath string
	header   string
	body     string
	formula  string
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
	for _, row := range body {
		b.body += b.separateData(row)
	}
}

func (b *CSVBuilder) build() {
	f, err := os.Create(b.filePath + ".csv")
	if err != nil {
		log.Fatal("Fail to build csv file: ", err)
	}
	defer f.Close()
	blankRow := b.separateData(make([]string, 0))

	fmt.Fprintf(f,
		b.header+
			blankRow+
			b.formula+
			blankRow+
			b.body,
	)
}

type formula int

const (
	fails formula = iota
	passes
	runs
	rate
)

var formulaNames = map[formula]string{
	fails:  "fails",
	passes: "passes",
	runs:   "runs",
	rate:   "rate",
}

func (f formula) name() string {
	return formulaNames[f]
}

// The known row IDs, prefix: [rid] as in [row-id].
type knownRowID int

const (
	_            knownRowID = iota // 0
	_                              // 1
	_                              // 2
	_                              // 3
	ridFails                       // 4
	ridPasses                      // 5
	ridRuns                        // 6
	ridPassRate                    // 7
	_                              // 8
	ridDataStart                   // 9
)

func (b *CSVBuilder) setFormula(totalRows int, totalCols int) {
	rowStart := int(ridDataStart)
	rowEnd := rowStart + totalRows - 1

	for _, f := range []formula{fails, passes, runs, rate} {
		var cells []string
		for col := 1; col <= totalCols; col++ {
			if col == 1 {
				cells = append(cells, f.name())
				continue
			}

			cellRange := fmt.Sprintf("%s:%s", NewTableCell(rowStart, col).id(), NewTableCell(rowEnd, col).id())
			switch f {
			case fails:
				cells = append(cells, fmt.Sprintf(`=COUNTIF(%s,"x")`, cellRange))
			case passes:
				cells = append(cells, fmt.Sprintf("=SUM(%s)", cellRange))
			case runs:
				cells = append(cells, fmt.Sprintf("=COUNTA(%s)", cellRange))
			case rate:
				cells = append(cells,
					fmt.Sprintf(`=IF(%[1]s=0,"N/A",%[2]s/%[1]s)`,
						NewTableCell(int(ridRuns), col).id(),
						NewTableCell(int(ridPasses), col).id(),
					))
			}
		}

		b.formula += b.separateData(cells)
	}
}

func (b *CSVBuilder) separateData(data []string) string {
	return fmt.Sprintf("%s\n", strings.Join(data, separator))
}

type TableCell struct {
	row int
	col int
}

func (cell *TableCell) id() string {
	return fmt.Sprintf("%s%d", cell.colNumToID(), cell.row)
}

func NewTableCell(row int, col int) *TableCell {
	return &TableCell{
		row: row,
		col: col,
	}
}

func (cell *TableCell) colNumToID() string {
	var res string
	//Convert columns to cell IDs (from left to right)
	for cell.col > 0 {
		mod := (cell.col - 1) % alphabetLength
		res = string(rune('A'+mod)) + res
		cell.col = (cell.col - mod) / alphabetLength
	}
	return res
}
