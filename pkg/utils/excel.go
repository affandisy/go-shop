package utils

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

type ExcelGenerator struct {
	file      *excelize.File
	sheetName string
	row       int
}

func NewExcelGenerator(sheetName string) *ExcelGenerator {
	f := excelize.NewFile()

	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	return &ExcelGenerator{
		file:      f,
		sheetName: sheetName,
		row:       1,
	}
}

func (g *ExcelGenerator) SetTitle(title string) {
	// Merge cells untuk title
	g.file.MergeCell(g.sheetName, "A1", "F1")

	// Set title
	g.file.SetCellValue(g.sheetName, "A1", title)

	// Style title
	style, _ := g.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Size:   16,
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	g.file.SetCellStyle(g.sheetName, "A1", "F1", style)
	g.file.SetRowHeight(g.sheetName, 1, 30)

	g.row = 3 // Skip to row 3
}

func (g *ExcelGenerator) AddTableHeader(headers []string) {
	// Header style
	style, _ := g.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:   true,
			Color:  "FFFFFF",
			Family: "Arial",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, header := range headers {
		cell := fmt.Sprintf("%c%d", 'A'+i, g.row)
		g.file.SetCellValue(g.sheetName, cell, header)
		g.file.SetCellStyle(g.sheetName, cell, cell, style)
	}

	g.row++
}

func (g *ExcelGenerator) AddTableRow(values []interface{}) {
	// Row style with borders
	style, _ := g.file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, value := range values {
		cell := fmt.Sprintf("%c%d", 'A'+i, g.row)
		g.file.SetCellValue(g.sheetName, cell, value)
		g.file.SetCellStyle(g.sheetName, cell, cell, style)
	}

	g.row++
}

func (g *ExcelGenerator) AddSummary(items map[string]interface{}) {
	g.row += 2 // Skip 2 rows

	// Summary style
	keyStyle, _ := g.file.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
	})

	for key, value := range items {
		keyCell := fmt.Sprintf("A%d", g.row)
		valueCell := fmt.Sprintf("B%d", g.row)

		g.file.SetCellValue(g.sheetName, keyCell, key)
		g.file.SetCellValue(g.sheetName, valueCell, value)
		g.file.SetCellStyle(g.sheetName, keyCell, keyCell, keyStyle)

		g.row++
	}
}

func (g *ExcelGenerator) AutoFitColumns(cols int) {
	for i := 0; i < cols; i++ {
		col := string(rune('A' + i))
		g.file.SetColWidth(g.sheetName, col, col, 20)
	}
}

func (g *ExcelGenerator) Output() ([]byte, error) {
	buf, err := g.file.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *ExcelGenerator) SaveToFile(filename string) error {
	return g.file.SaveAs(filename)
}
