package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func buildItemTable(state *UIState) fyne.CanvasObject {
	state.Items = container.NewVBox()

	header := container.NewGridWithColumns(5,
		widget.NewLabelWithStyle("Item", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Quantity", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Rate", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Amount", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
	)

	addItemBtn := widget.NewButtonWithIcon("Line Item", theme.ContentAddIcon(), func() {
		state.AddItemRow()
		state.CalculateTotals()
	})
	addItemBtn.Importance = widget.HighImportance

	// Add initial row
	state.AddItemRow()

	return container.NewVBox(header, state.Items, addItemBtn)
}

func (s *UIState) AddItemRow() {
	desc := widget.NewEntry()
	desc.SetPlaceHolder("Description of item/service...")

	qty := widget.NewEntry()
	qty.SetText("1")
	qty.OnChanged = func(string) { s.CalculateTotals() }

	rate := widget.NewEntry()
	rate.SetText("0")
	rate.OnChanged = func(string) { s.CalculateTotals() }

	amount := widget.NewLabel("$0.00")

	removeBtn := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)

	rowLayout := container.NewGridWithColumns(5,
		desc,
		qty,
		rate,
		amount,
		removeBtn,
	)

	row := &ItemRow{
		Description: desc,
		Quantity:    qty,
		Rate:        rate,
		AmountLabel: amount,
		Container:   rowLayout,
	}

	removeBtn.OnTapped = func() {
		s.Items.Remove(rowLayout)
		for i, r := range s.ItemRows {
			if r == row {
				s.ItemRows = append(s.ItemRows[:i], s.ItemRows[i+1:]...)
				break
			}
		}
		s.CalculateTotals()
	}

	s.ItemRows = append(s.ItemRows, row)
	s.Items.Add(rowLayout)
}
