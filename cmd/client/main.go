package main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"licensebox/internal/client"
	invoice_ui "licensebox/internal/client/ui"
)

type LicenseUI struct {
	App     fyne.App
	Window  fyne.Window
	Manager *client.LicenseManager
}

func main() {
	myApp := app.NewWithID("com.licensebox.client")
	myApp.Settings().SetTheme(theme.DarkTheme())

	myWindow := myApp.NewWindow("Software Activation")
	myWindow.Resize(fyne.NewSize(500, 350))
	myWindow.CenterOnScreen()

	manager := client.NewLicenseManager("http://localhost:8080")
	u := &LicenseUI{App: myApp, Window: myWindow, Manager: manager}

	u.showLoading()
	go u.checkLicenseStatus()

	myWindow.ShowAndRun()
}

func (ui *LicenseUI) showLoading() {
	loading := widget.NewProgressBarInfinite()
	label := widget.NewLabelWithStyle("Checking License...", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	content := container.NewVBox(
		layout.NewSpacer(),
		label,
		loading,
		layout.NewSpacer(),
	)
	ui.Window.SetContent(content)
}

func (ui *LicenseUI) checkLicenseStatus() {
	active, err := ui.Manager.CheckLicense()
	if err != nil {
		fmt.Printf("Error checking license: %v\n", err)
	}

	// Update UI on main thread
	fyne.Do(func() {
		if active {
			ui.showMainApp()
		} else {
			ui.showActivationScreen()
		}
	})
}

func (ui *LicenseUI) showActivationScreen() {
	title := canvas.NewText("Activate Your Software", color.NRGBA{R: 0, G: 200, B: 200, A: 255})
	title.TextSize = 24
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	subtitle := widget.NewLabel("Please enter your license key to continue")
	subtitle.Alignment = fyne.TextAlignCenter

	licenseEntry := widget.NewEntry()
	licenseEntry.SetPlaceHolder("XXXX-XXXX-XXXX-XXXX")

	activateBtn := widget.NewButtonWithIcon("Activate Now", theme.ConfirmIcon(), func() {
		key := licenseEntry.Text
		if key == "" {
			dialog.ShowError(fmt.Errorf("License key is required"), ui.Window)
			return
		}

		progressContent := container.NewVBox(widget.NewLabel("Connecting to server..."), widget.NewProgressBarInfinite())
		progress := dialog.NewCustomWithoutButtons("Activating", progressContent, ui.Window)
		progress.Show()

		go func() {
			err := ui.Manager.Activate(key)
			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					dialog.ShowError(err, ui.Window)
				} else {
					dialog.ShowInformation("Success", "Software activated successfully!", ui.Window)
					ui.showMainApp()
				}
			})
		}()
	})
	activateBtn.Importance = widget.HighImportance

	content := container.New(layout.NewVBoxLayout(),
		container.NewCenter(title),
		container.NewCenter(subtitle),
		layout.NewSpacer(),
		licenseEntry,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), activateBtn, layout.NewSpacer()),
		layout.NewSpacer(),
	)

	paddedContent := container.NewPadded(content)
	ui.Window.SetContent(paddedContent)
}

func (ui *LicenseUI) showMainApp() {
	content := invoice_ui.BuildInvoiceUI(ui.Window, ui.App)
	ui.Window.SetContent(content)
}
