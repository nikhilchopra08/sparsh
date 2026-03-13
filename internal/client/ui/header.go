package ui

import (
	"image/color"
	"licensebox/internal/client"
	"licensebox/internal/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	dialog_fyne "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/sqweek/dialog"
)

func buildHeader(state *UIState) fyne.CanvasObject {
	state.LogoImage = &canvas.Image{FillMode: canvas.ImageFillContain}
	state.LogoImage.SetMinSize(fyne.NewSize(120, 80))

	pickLogo := func() {
		filename, err := dialog.File().Title("Select Logo").Filter("Image Files", "png", "jpg", "jpeg").Load()
		if err != nil {
			if err.Error() != "Cancelled" {
				dialog_fyne.ShowError(err, state.Window)
			}
			return
		}
		state.LogoPath = filename
		state.LogoImage.File = filename
		state.LogoImage.Refresh()
		state.LogoContainer.Show()
	}

	logoBtn := widget.NewButton("+ Add Your Logo", pickLogo)

	clickableLogo := widget.NewButton("", pickLogo)
	clickableLogo.Importance = widget.LowImportance

	state.LogoContainer = container.NewStack(state.LogoImage, clickableLogo)
	state.LogoContainer.Hide()

	if settings, _ := client.LoadUserSettings(); settings != nil && settings.LogoPath != "" {
		state.LogoPath = settings.LogoPath
		state.LogoImage.File = settings.LogoPath
		state.LogoContainer.Show()
	}

	logoStack := container.NewStack(logoBtn, state.LogoContainer)
	logoBox := container.NewPadded(logoStack)

	// Action buttons rows
	row1Btns := container.NewHBox(
		widget.NewButton("Save Template", func() {
			invoice := state.GatherTemplateData()
			client.SaveInvoice(invoice, "template.json")
			dialog_fyne.ShowInformation("Saved", "Template saved.", state.Window)
		}),
		widget.NewButton("Load", func() {
			inv, _ := client.LoadInvoice("template.json")
			if inv != nil {
				state.SetData(*inv)
			}
		}),
		widget.NewButton("Export JSON", func() {
			invoice := state.GatherData()
			filename, err := dialog.File().Title("Export JSON").Filter("JSON", "json").Save()
			if err == nil && filename != "" {
				client.SaveInvoice(invoice, filename)
			}
		}),
	)

	row2Btns := container.NewHBox(
		widget.NewButton("Import JSON", func() {
			filename, err := dialog.File().Title("Import JSON").Filter("JSON", "json").Load()
			if err == nil && filename != "" {
				inv, _ := client.LoadInvoice(filename)
				if inv != nil {
					state.SetData(*inv)
				}
			}
		}),
	)

	actionButtons := container.NewVBox(row1Btns, row2Btns)
	leftSide := container.NewHBox(logoBox, actionButtons)

	// -- Right Side --
	invoiceHeader := canvas.NewText("INVOICE", color.NRGBA{R: 28, G: 40, B: 51, A: 255})
	invoiceHeader.TextSize = 36
	invoiceHeader.TextStyle = fyne.TextStyle{Bold: true}

	// Brand section
	brandLabel := canvas.NewText("Brand", color.NRGBA{R: 128, G: 139, B: 150, A: 255})
	brandLabel.TextSize = 11
	brandSelect := widget.NewSelect([]string{"Default", "Brand A", "Brand B"}, func(s string) {})
	brandSelect.SetSelected("Default")

	manageBtn := widget.NewButton("Manage", func() {})
	brandWithLabel := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), brandLabel, layout.NewSpacer()),
		manageBtn,
	)

	// Invoice Number section
	hashLabel := canvas.NewText("#", color.NRGBA{R: 128, G: 139, B: 150, A: 255})
	hashLabel.TextSize = 11
	state.InvoiceNum = widget.NewEntry()
	state.InvoiceNum.SetText("001")
	if settings, _ := client.LoadUserSettings(); settings != nil && settings.LastInvoiceNum != "" {
		state.InvoiceNum.SetText(settings.LastInvoiceNum)
	}
	invoiceNumWithLabel := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), hashLabel, layout.NewSpacer()),
		state.InvoiceNum,
	)

	// Next # block
	nextText := canvas.NewText("Next", color.NRGBA{R: 28, G: 40, B: 51, A: 255})
	nextText.TextSize = 12
	hashText := canvas.NewText("#", color.NRGBA{R: 28, G: 40, B: 51, A: 255})
	hashText.TextSize = 12
	nextBlock := container.NewVBox(nextText, hashText)
	nextBtn := widget.NewButton("", func() {
		state.IncrementInvoiceNumber()
	})
	nextSquare := container.NewStack(nextBtn, container.NewCenter(nextBlock))

	controlsRow := container.NewHBox(
		layout.NewSpacer(),
		brandSelect,
		brandWithLabel,
		invoiceNumWithLabel,
		nextSquare,
	)

	headerRight := container.NewVBox(
		container.NewHBox(layout.NewSpacer(), invoiceHeader),
		controlsRow,
	)

	return container.NewBorder(nil, nil, leftSide, headerRight)
}

func buildSettingsBar(state *UIState) fyne.CanvasObject {
	themeSelect := widget.NewSelect([]string{"Classic", "Modern"}, func(s string) {})
	themeSelect.SetSelected("Classic")
	currencySelect := widget.NewSelect([]string{"USD ($)", "EUR (€)", "GBP (£)"}, func(s string) {})
	currencySelect.SetSelected("USD ($)")

	state.AutoIncrement = widget.NewCheck("Auto-increment on PDF", func(b bool) {})
	state.AutoIncrement.SetChecked(true)

	saveDefaultBtn := widget.NewButton("Save Default", func() {
		client.SaveUserSettings(models.UserSettings{
			CompanyName:    state.FromEntry.Text,
			LastInvoiceNum: state.InvoiceNum.Text,
			LogoPath:       state.LogoPath,
		})
		dialog_fyne.ShowInformation("Saved", "Settings saved.", state.Window)
	})

	return container.NewHBox(
		widget.NewLabel("Theme"), themeSelect,
		widget.NewLabel("Currency"), currencySelect,
		widget.NewCheck("Show converted total", func(b bool) {}),
		layout.NewSpacer(),
		saveDefaultBtn,
		state.AutoIncrement,
	)
}
