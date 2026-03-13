package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func buildTotalsSection(state *UIState) fyne.CanvasObject {
	state.Notes = widget.NewMultiLineEntry()
	state.Notes.SetPlaceHolder("Any relevant information not already covered")

	state.PolicyTerms = widget.NewMultiLineEntry()
	state.PolicyTerms.SetPlaceHolder("Late fees, payment methods, delivery schedule")

	state.TaxRate = widget.NewEntry()
	state.TaxRate.SetText("0")
	state.TaxRate.OnChanged = func(string) { state.CalculateTotals() }

	state.AmountPaid = widget.NewEntry()
	state.AmountPaid.SetText("0")
	state.AmountPaid.OnChanged = func(string) { state.CalculateTotals() }

	state.SubtotalLabel = widget.NewLabel("$0.00")
	state.TotalLabel = widget.NewLabelWithStyle("$0.00", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})
	state.BalanceDueLabel = widget.NewLabelWithStyle("$0.00", fyne.TextAlignTrailing, fyne.TextStyle{Bold: true})

	totals := container.NewVBox(
		container.NewHBox(widget.NewLabel("Subtotal"), layout.NewSpacer(), state.SubtotalLabel),
		container.NewHBox(widget.NewLabel("Tax Rate (%)"), state.TaxRate, layout.NewSpacer(), widget.NewLabel("$0.00")),
		container.NewHBox(widget.NewLabelWithStyle("Total", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), state.TotalLabel),
		container.NewHBox(widget.NewLabel("Amount Paid"), state.AmountPaid),
		container.NewHBox(widget.NewLabelWithStyle("Balance Due", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), layout.NewSpacer(), state.BalanceDueLabel),
	)

	return container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabelWithStyle("Notes", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), state.Notes,
			widget.NewLabelWithStyle("Terms", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), state.PolicyTerms,
		),
		container.NewPadded(totals),
	)
}
