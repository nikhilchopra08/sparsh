package ui

import (
	"fmt"
	"licensebox/internal/client"
	"licensebox/internal/models"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

type UIState struct {
	Window fyne.Window
	App    fyne.App

	LogoPath      string
	LogoImage     *canvas.Image
	LogoContainer *fyne.Container

	InvoiceNum    *widget.Entry
	FromEntry     *widget.Entry
	BillToName    *widget.Entry
	BillToAddr    *widget.Entry
	ShipToName    *widget.Entry
	ShipToAddr    *widget.Entry
	Date          *widget.Entry
	DueDate       *widget.Entry
	Terms         *widget.Entry
	PO            *widget.Entry
	Notes         *widget.Entry
	PolicyTerms   *widget.Entry
	AutoIncrement *widget.Check

	TaxRate    *widget.Entry
	AmountPaid *widget.Entry

	SubtotalLabel   *widget.Label
	TotalLabel      *widget.Label
	BalanceDueLabel *widget.Label

	Items    *fyne.Container
	ItemRows []*ItemRow
}

type ItemRow struct {
	Description *widget.Entry
	Quantity    *widget.Entry
	Rate        *widget.Entry
	AmountLabel *widget.Label
	Container   *fyne.Container
}

func (s *UIState) CalculateTotals() {
	var subtotal float64
	for _, row := range s.ItemRows {
		q := parseFloat(row.Quantity.Text)
		r := parseFloat(row.Rate.Text)
		amount := q * r
		row.AmountLabel.SetText(formatCurrency(amount))
		subtotal += amount
	}

	taxRate := parseFloat(s.TaxRate.Text)
	taxAmount := subtotal * (taxRate / 100.0)
	total := subtotal + taxAmount
	paid := parseFloat(s.AmountPaid.Text)
	balance := total - paid

	s.SubtotalLabel.SetText(formatCurrency(subtotal))
	s.TotalLabel.SetText(formatCurrency(total))
	s.BalanceDueLabel.SetText(formatCurrency(balance))
}

func (s *UIState) GatherData() models.Invoice {
	taxRate := parseFloat(s.TaxRate.Text)
	paid := parseFloat(s.AmountPaid.Text)

	items := make([]models.InvoiceItem, 0)
	for _, row := range s.ItemRows {
		items = append(items, models.InvoiceItem{
			Description: row.Description.Text,
			Quantity:    parseFloat(row.Quantity.Text),
			Rate:        parseFloat(row.Rate.Text),
			Amount:      parseFloat(row.Quantity.Text) * parseFloat(row.Rate.Text),
		})
	}

	return models.Invoice{
		Number:       s.InvoiceNum.Text,
		Date:         s.Date.Text,
		DueDate:      s.DueDate.Text,
		PaymentTerms: s.Terms.Text,
		PONumber:     s.PO.Text,
		From:         s.FromEntry.Text,
		BillToName:   s.BillToName.Text,
		BillToAddr:   s.BillToAddr.Text,
		ShipToName:   s.ShipToName.Text,
		ShipToAddr:   s.ShipToAddr.Text,
		Items:        items,
		Notes:        s.Notes.Text,
		Terms:        s.PolicyTerms.Text,
		TaxRate:      taxRate,
		AmountPaid:   paid,
		LogoPath:     s.LogoPath,
	}
}
func (s *UIState) GatherTemplateData() models.Invoice {
	taxRate := parseFloat(s.TaxRate.Text)
	paid := parseFloat(s.AmountPaid.Text)

	items := make([]models.InvoiceItem, 0)
	for _, row := range s.ItemRows {
		items = append(items, models.InvoiceItem{
			Description: "", // Clear description for templates
			Quantity:    parseFloat(row.Quantity.Text),
			Rate:        parseFloat(row.Rate.Text),
			Amount:      parseFloat(row.Quantity.Text) * parseFloat(row.Rate.Text),
		})
	}

	return models.Invoice{
		Number:       "", // Exclude invoice number from templates
		Date:         "", // User wants current date on load, so template leaves blank
		DueDate:      "",
		PaymentTerms: s.Terms.Text,
		PONumber:     "",
		From:         s.FromEntry.Text, // Preserve branding
		BillToName:   "",
		BillToAddr:   "",
		ShipToName:   "",
		ShipToAddr:   "",
		Items:        items,
		Notes:        s.Notes.Text,
		Terms:        s.PolicyTerms.Text,
		TaxRate:      taxRate,
		AmountPaid:   paid,
		LogoPath:     s.LogoPath,
	}
}

func (s *UIState) SetData(inv models.Invoice) {
	// Logo handling
	s.LogoPath = inv.LogoPath
	if s.LogoPath != "" {
		s.LogoImage.File = s.LogoPath
		s.LogoImage.Refresh()
		s.LogoContainer.Show()
	} else {
		s.LogoImage.File = ""
		s.LogoContainer.Hide()
	}

	if inv.Number != "" {
		s.InvoiceNum.SetText(inv.Number)
	}

	today := time.Now().Format("02/01/2006")
	if inv.Date != "" {
		s.Date.SetText(inv.Date)
	} else {
		s.Date.SetText(today)
	}

	if inv.DueDate != "" {
		s.DueDate.SetText(inv.DueDate)
	} else {
		s.DueDate.SetText(today)
	}

	s.Terms.SetText(inv.PaymentTerms)
	s.PO.SetText(inv.PONumber)
	s.FromEntry.SetText(inv.From)
	s.BillToName.SetText(inv.BillToName)
	s.BillToAddr.SetText(inv.BillToAddr)
	s.ShipToName.SetText(inv.ShipToName)
	s.ShipToAddr.SetText(inv.ShipToAddr)
	s.Notes.SetText(inv.Notes)
	s.PolicyTerms.SetText(inv.Terms)
	s.TaxRate.SetText(fmt.Sprintf("%.1f", inv.TaxRate))
	s.AmountPaid.SetText(fmt.Sprintf("%.2f", inv.AmountPaid))

	// Rebuild items
	s.Items.Objects = nil
	s.ItemRows = nil
	for _, item := range inv.Items {
		s.AddItemRow()
		row := s.ItemRows[len(s.ItemRows)-1]
		row.Description.SetText(item.Description)
		row.Quantity.SetText(fmt.Sprintf("%.1f", item.Quantity))
		row.Rate.SetText(fmt.Sprintf("%.2f", item.Rate))
	}
	s.CalculateTotals()
}

func (s *UIState) IncrementInvoiceNumber() {
	numStr := s.InvoiceNum.Text
	if numStr == "" {
		s.InvoiceNum.SetText("001")
		return
	}

	// Try to find numeric part
	var num int
	fmt.Sscanf(numStr, "%d", &num)

	newNum := num + 1
	// Simple padding to 3 digits
	nextStr := fmt.Sprintf("%03d", newNum)
	s.InvoiceNum.SetText(nextStr)

	// Save to persistent settings
	settings, _ := client.LoadUserSettings()
	if settings == nil {
		settings = &models.UserSettings{}
	}
	settings.LastInvoiceNum = nextStr
	client.SaveUserSettings(*settings)
}

func parseFloat(s string) float64 {
	var val float64
	fmt.Sscanf(s, "%f", &val)
	return val
}

func formatCurrency(val float64) string {
	return fmt.Sprintf("$%.2f", val)
}
