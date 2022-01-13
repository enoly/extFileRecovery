package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ext_worker "github.com/enoly/extFileRecovery/internal/extworker"
)

func getRestoreIndirectCallback(name string, inode *[]uint32, worker *ext_worker.Ext3Worker) func() {
	n := name
	i := *inode
	return func() {
		worker.RestoreFromPtrs(n, &i)
	}
}

func showFoundIndirect(dir string, worker *ext_worker.Ext3Worker, window fyne.Window) {
	found, err := worker.RestoreFromIndirectBlocks(dir)
	if err != nil {
		errorLabel := widget.NewLabel(fmt.Sprintf("Unable to find deleted from directory: %v in indirect\nerror: %v", dir, err))
		backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
		content := container.New(layout.NewVBoxLayout(), errorLabel, backButton)
		window.SetContent(container.New(layout.NewCenterLayout(), content))
		return
	}

	if len(*found) == 0 {
		errorLabel := widget.NewLabel(fmt.Sprintf("Found nothing to restore from directory: %v", dir))
		backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
		content := container.New(layout.NewVBoxLayout(), errorLabel, backButton)
		window.SetContent(container.New(layout.NewCenterLayout(), content))
		return
	}

	form := container.New(layout.NewFormLayout())
	for name, arr := range *found {
		for i, ptrs := range arr {
			filePath := ""
			if dir == "/" {
				filePath = "/" + name
			} else if dir[len(dir)-1] == '/' {
				filePath = dir[:len(dir)-2] + name
			} else {
				filePath = dir + string('/') + name
			}
			foundLabel := widget.NewLabel(fmt.Sprintf("file: %v\nsize: %vb", filePath, len(ptrs)*int(worker.ExtFs.BlockSize)))
			restoreButton := widget.NewButton("Restore", getRestoreIndirectCallback(fmt.Sprintf("%v%v", i, name), &ptrs, worker))
			form.Add(foundLabel)
			form.Add(restoreButton)
		}
	}
	header := widget.NewLabel(fmt.Sprintf("Found %v inodes to restore from indirect", len(*found)))
	backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
	content := container.New(layout.NewVBoxLayout(), header, form, backButton)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func OpenFindIndirect(dir string, worker *ext_worker.Ext3Worker, window fyne.Window) {
	showSearch(window)
	go showFoundIndirect(dir, worker, window)
}
