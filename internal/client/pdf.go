package client

import (
	"fmt"
	"licensebox/internal/models"

	"github.com/jung-kurt/gofpdf"
)

func ExportInvoiceToPDF(invoice models.Invoice, outputPath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Logo
	if invoice.LogoPath != "" {
		pdf.ImageOptions(invoice.LogoPath, 10, 10, 0, 20, false, gofpdf.ImageOptions{ImageType: "", ReadDpi: true}, 0, "")
		pdf.Ln(25) // Space for logo
	} else {
		// Header
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, "INVOICE")
		pdf.Ln(12)
	}

	if invoice.LogoPath != "" {
		pdf.SetFont("Arial", "B", 16)
		pdf.Cell(40, 10, "INVOICE")
		pdf.Ln(12)
	}

	// Invoice Info
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 10, fmt.Sprintf("Invoice #: %s", invoice.Number))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Date: %s", invoice.Date))
	pdf.Ln(6)
	pdf.Cell(0, 10, fmt.Sprintf("Due Date: %s", invoice.DueDate))
	pdf.Ln(10)

	// Addresses
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(90, 10, "From:")
	pdf.Cell(90, 10, "Bill To:")
	pdf.Ln(6)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(90, 10, invoice.From)
	pdf.Cell(90, 10, fmt.Sprintf("%s\n%s", invoice.BillToName, invoice.BillToAddr))
	pdf.Ln(20)

	// Items Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(100, 10, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 10, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Rate", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 10, "Amount", "1", 1, "C", true, 0, "")

	// Items
	pdf.SetFont("Arial", "", 10)
	var subtotal float64
	for _, item := range invoice.Items {
		pdf.CellFormat(100, 8, item.Description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", item.Rate), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", item.Amount), "1", 1, "C", false, 0, "")
		subtotal += item.Amount
	}

	// Totals
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(160, 8, "Subtotal:")
	pdf.Cell(30, 8, fmt.Sprintf("%.2f", subtotal))
	pdf.Ln(6)

	taxAmount := subtotal * (invoice.TaxRate / 100)
	pdf.Cell(160, 8, fmt.Sprintf("Tax (%.1f%%):", invoice.TaxRate))
	pdf.Cell(30, 8, fmt.Sprintf("%.2f", taxAmount))
	pdf.Ln(6)

	total := subtotal + taxAmount
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(160, 10, "Total:")
	pdf.Cell(30, 10, fmt.Sprintf("%.2f", total))
	pdf.Ln(10)

	// Notes
	if invoice.Notes != "" {
		pdf.SetFont("Arial", "B", 10)
		pdf.Cell(0, 10, "Notes:")
		pdf.Ln(6)
		pdf.SetFont("Arial", "", 10)
		pdf.MultiCell(0, 5, invoice.Notes, "", "L", false)
	}

	return pdf.OutputFileAndClose(outputPath)
}
