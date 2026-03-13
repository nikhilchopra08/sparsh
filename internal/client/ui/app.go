package ui

import (
	"fmt"
	"image/color"
	"licensebox/internal/client"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	dialog_fyne "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func BuildInvoiceUI(window fyne.Window, app fyne.App) fyne.CanvasObject {
	app.Settings().SetTheme(theme.LightTheme())
	window.Resize(fyne.NewSize(1000, 800))

	state := &UIState{
		Window:   window,
		App:      app,
		ItemRows: make([]*ItemRow, 0),
	}

	// 1. Initialize Components
	header := buildHeader(state)
	settings := buildSettingsBar(state)
	address := buildAddressSection(state)
	details := buildDetailsGrid(state)
	items := buildItemTable(state)
	totals := buildTotalsSection(state)

	downloadBtn := widget.NewButtonWithIcon("Download PDF", theme.DownloadIcon(), func() {
		invoice := state.GatherData()
		defaultName := "invoice.pdf"
		if invoice.Number != "" {
			defaultName = fmt.Sprintf("invoice_%s.pdf", invoice.Number)
		}

		// Native Save Dialog
		filename, err := dialog.File().Title("Save Invoice PDF").Filter("PDF Files", "pdf").SetStartFile(defaultName).Save()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog_fyne.ShowError(err, window)
			}
			return
		}

		err = client.ExportInvoiceToPDF(invoice, filename)
		if err != nil {
			dialog_fyne.ShowError(err, window)
		} else {
			if state.AutoIncrement.Checked {
				state.IncrementInvoiceNumber()
			}
			dialog_fyne.ShowInformation("Export Success", fmt.Sprintf("Invoice downloaded to %s", filename), window)
		}
	})
	downloadBtn.Importance = widget.HighImportance

	content := container.NewVBox(
		header,
		settings,
		canvas.NewLine(color.Gray{Y: 0xcc}),
		address,
		details,
		items,
		totals,
	)

	scroll := container.NewVScroll(container.NewPadded(content))

	return container.NewBorder(nil, container.NewHBox(layout.NewSpacer(), downloadBtn), nil, nil, scroll)
}
