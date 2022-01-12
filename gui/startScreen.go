package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func OpenStartScreen(window fyne.Window) {
	header := widget.NewLabel("extFileRecovery")
	description := widget.NewLabel("ext3-4 file recovery tool")
	description.Alignment = fyne.TextAlignCenter
	header.Alignment = fyne.TextAlignCenter

	selectButton := widget.NewButton("Select drive", func() {
		OpenSelectDrive(window)
	})

	innerContent := container.New(layout.NewVBoxLayout(), header, description, selectButton)
	window.SetContent(container.New(layout.NewCenterLayout(), innerContent))
}
