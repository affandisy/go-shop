package utils

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type PDFGenerator struct {
	pdf *gofpdf.Fpdf
}

func NewPDFGenerator() *PDFGenerator {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 15)

	return &PDFGenerator{pdf: pdf}
}

func (g *PDFGenerator) SetTitle(title string) {
	g.pdf.AddPage()

	g.pdf.SetFont("Arial", "B", 20)
	g.pdf.CellFormat(0, 15, title, "", 1, "C", false, 0, "")

	g.pdf.SetFont("Arial", "", 10)
	dateStr := fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05"))
	g.pdf.CellFormat(0, 5, dateStr, "", 1, "C", false, 0, "")
	g.pdf.Ln(5)
}

func (g *PDFGenerator) AddTableHeader(headers []string, widths []float64) {
	g.pdf.SetFont("Arial", "B", 10)
	g.pdf.SetFillColor(200, 220, 255)

	for i, header := range headers {
		g.pdf.CellFormat(widths[i], 8, header, "1", 0, "C", true, 0, "")
	}
	g.pdf.Ln(-1)
}

func (g *PDFGenerator) AddTableRow(values []string, widths []float64) {
	g.pdf.SetFont("Arial", "", 9)

	for i, value := range values {
		g.pdf.CellFormat(widths[i], 7, value, "1", 0, "L", false, 0, "")
	}

	g.pdf.Ln(-1)
}

func (g *PDFGenerator) AddSummary(title string, items map[string]string) {
	g.pdf.Ln(5)
	g.pdf.SetFont("Arial", "B", 12)
	g.pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")

	g.pdf.SetFont("Arial", "", 10)
	for key, value := range items {
		g.pdf.CellFormat(70, 6, key, "", 0, "L", false, 0, "")
		g.pdf.CellFormat(0, 6, ": "+value, "", 1, "L", false, 0, "")
	}
}

func (g *PDFGenerator) AddText(text string) {
	g.pdf.SetFont("Arial", "", 10)
	g.pdf.MultiCell(0, 5, text, "", "L", false)
	g.pdf.Ln(3)
}

func (g *PDFGenerator) Output() ([]byte, error) {
	var buf []byte
	writer := &bytesWriter{buf: &buf}

	err := g.pdf.Output(writer)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (g *PDFGenerator) SaveToFile(filename string) error {
	return g.pdf.OutputFileAndClose(filename)
}

type bytesWriter struct {
	buf *[]byte
}

func (w *bytesWriter) Write(p []byte) (n int, err error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}
