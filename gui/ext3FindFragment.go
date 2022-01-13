package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ext_worker "github.com/enoly/extFileRecovery/internal/extworker"
)

func showFinish(window fyne.Window) {
	finishLabel := widget.NewLabel("Search completed.\nRestored files are in programm directory.")
	backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
	content := container.New(layout.NewVBoxLayout(), finishLabel, backButton)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func findFragment(worker *ext_worker.Ext3Worker, counter *binding.Float, window fyne.Window) {
	found := make(chan uint64, 10)

	go worker.RestoreFragments(found, counter)
	for {
		_, ok := <-found
		if !ok {
			break
		}
	}

	showFinish(window)
}

func OpenFindFragment(worker *ext_worker.Ext3Worker, window fyne.Window) {
	counter := binding.NewFloat()
	counter.Set(0)
	progress := widget.NewProgressBarWithData(counter)
	findText := widget.NewLabel("Search in progress...")
	content := container.New(layout.NewVBoxLayout(), findText, progress)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
	go findFragment(worker, &counter, window)
}
