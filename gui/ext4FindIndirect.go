package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ext_worker "github.com/enoly/extFileRecovery/internal/extworker"
	"github.com/enoly/extFileRecovery/pkg/ext4/structure"
)

func getRestoreFunc(extent *structure.Extent, i int, worker *ext_worker.Ext4Worker, window fyne.Window) func() {
	var e structure.Extent = *extent
	return func() {
		saved := window.Content()
		window.SetContent(container.New(layout.NewCenterLayout(), widget.NewLabel("Restoring...")))
		if err := worker.RestoreFileFromExtent(&e, i); err != nil {
			panic(err)
		}
		window.SetContent(saved)
	}
}

func openFindingResults(extents *[]uint64, worker *ext_worker.Ext4Worker, window fyne.Window) {
	if len(*extents) == 0 {
		notFoundText := widget.NewLabel("Found 0 indirect extents")
		backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
		content := container.New(layout.NewVBoxLayout(), notFoundText, backButton)
		window.SetContent(container.New(layout.NewCenterLayout(), content))
		return
	}

	form := container.New(layout.NewFormLayout())
	for _, exextentBlock := range *extents {
		extent, err := (*worker).GetExtentFromBlock(exextentBlock)
		if err != nil {
			continue
		}

		if len(extent.LeafNodes) == 0 {
			continue
		}

		fileSize := 0
		for _, leaf := range extent.LeafNodes {
			fileSize += int(leaf.CoveredBlocks) * int(worker.ExtFs.BlockSize)
		}
		extentText := fmt.Sprintf("Block: %v\nFile size: %vb", exextentBlock, fileSize)
		extentLabel := widget.NewLabel(extentText)
		restoreButton := widget.NewButton("Restore", getRestoreFunc(extent, int(exextentBlock), worker, window))
		form.Add(extentLabel)
		form.Add(restoreButton)
	}

	header := widget.NewLabel(fmt.Sprintf("Found %v indirect extents", len(*extents)))
	backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
	content := container.New(layout.NewVBoxLayout(), header, form, backButton)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func findExt4Indirect(worker *ext_worker.Ext4Worker, counter *binding.Float, window fyne.Window) {
	found := make(chan uint64, 10)
	extents := []uint64{}

	go worker.FindIndirectExtents(found, counter)
	for {
		f, ok := <-found
		if !ok {
			break
		}

		extents = append(extents, f)
	}

	openFindingResults(&extents, worker, window)
}

func OpenExt4FindIndirect(worker *ext_worker.Ext4Worker, window fyne.Window) {
	counter := binding.NewFloat()
	counter.Set(0)
	progress := widget.NewProgressBarWithData(counter)
	findText := widget.NewLabel("Search in progress...")
	content := container.New(layout.NewVBoxLayout(), findText, progress)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
	go findExt4Indirect(worker, &counter, window)
}
