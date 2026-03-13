package ui

import (
	"licensebox/internal/client"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func buildAddressSection(state *UIState) fyne.CanvasObject {
	state.FromEntry = widget.NewMultiLineEntry()
	state.FromEntry.SetPlaceHolder("Your company name, address, etc.")

	// Pre-load from settings
	settings, _ := client.LoadUserSettings()
	if settings != nil && settings.CompanyName != "" {
		state.FromEntry.SetText(settings.CompanyName)
	}

	state.BillToName = widget.NewEntry()
	state.BillToName.SetPlaceHolder("Client name")
	state.BillToAddr = widget.NewMultiLineEntry()
	state.BillToAddr.SetPlaceHolder("Client address")

	state.ShipToName = widget.NewEntry()
	state.ShipToName.SetPlaceHolder("Recipient / Department")
	state.ShipToAddr = widget.NewMultiLineEntry()
	state.ShipToAddr.SetPlaceHolder("Shipping address")

	return container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabelWithStyle("Who is this from?", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), state.FromEntry),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabelWithStyle("Bill To", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), state.BillToName, state.BillToAddr),
			container.NewVBox(widget.NewLabelWithStyle("Ship To (optional)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), state.ShipToName, state.ShipToAddr),
		),
	)
}

func buildDetailsGrid(state *UIState) fyne.CanvasObject {
	today := time.Now().Format("02/01/2006")
	state.Date = widget.NewEntry()
	state.Date.SetText(today)
	state.DueDate = widget.NewEntry()
	state.DueDate.SetText(today)
	state.Terms = widget.NewEntry()
	state.Terms.SetPlaceHolder("e.g. Net 15")
	state.PO = widget.NewEntry()
	state.PO.SetPlaceHolder("PO Number")

	return container.NewGridWithColumns(2,
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel("Date"), state.Date),
			container.NewVBox(widget.NewLabel("Payment Terms"), state.Terms),
		),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel("Due Date"), state.DueDate),
			container.NewVBox(widget.NewLabel("PO Number"), state.PO),
		),
	)
}
