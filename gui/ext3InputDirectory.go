package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ext_worker "github.com/enoly/extFileRecovery/internal/extworker"
)

func showError(err error, worker *ext_worker.Ext3Worker, window fyne.Window) {
	content := getInputContent(worker, window)
	content.Add(widget.NewLabel(fmt.Sprintf("Unable to select directory:\n%v", err)))
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func showButtons(dir string, worker *ext_worker.Ext3Worker, window fyne.Window) {
	journal := widget.NewButton("Find in journal", func() { OpenFindJournal(dir, worker, window) })
	indirect := widget.NewButton("Find indirect", func() { OpenFindIndirect(dir, worker, window) })

	content := container.New(layout.NewVBoxLayout(), journal, indirect)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func getInputContent(worker *ext_worker.Ext3Worker, window fyne.Window) *fyne.Container {
	header := widget.NewLabel("Please input directory for search")
	header.Alignment = fyne.TextAlignCenter

	input := widget.NewEntry()
	input.SetText("/")

	button := widget.NewButton("Next", func() {
		_, err := worker.ReadDirectory(nil, input.Text)
		if input.Text == "" {
			input.SetText("/")
		}

		if err != nil {
			showError(err, worker, window)
		} else {
			showButtons(input.Text, worker, window)
		}
	})

	return container.New(layout.NewVBoxLayout(), header, input, button)
}

func OpenInputDirectory(worker *ext_worker.Ext3Worker, window fyne.Window) {
	window.SetContent(container.New(layout.NewCenterLayout(), getInputContent(worker, window)))
}
