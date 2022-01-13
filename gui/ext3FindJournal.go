package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	ext_worker "github.com/enoly/extFileRecovery/internal/extworker"
	structure "github.com/enoly/extFileRecovery/pkg/ext3/structure"
)

func showSearch(window fyne.Window) {
	label := widget.NewLabel("Search in progress...")
	window.SetContent(container.New(layout.NewCenterLayout(), label))
}

func getRestoreCallback(name string, inode *structure.Inode, worker *ext_worker.Ext3Worker) func() {
	n := name
	i := *inode
	return func() {
		worker.RestoreFromInode(n, &i)
	}
}

func showFound(dir string, worker *ext_worker.Ext3Worker, window fyne.Window) {
	found, err := worker.FindInJournal(dir)
	if err != nil || len(*found) == 0 {
		errorLabel := widget.NewLabel(fmt.Sprintf("Unable to find deleted from directory: %v in journal\nerror: %v", dir, err))
		backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
		content := container.New(layout.NewVBoxLayout(), errorLabel, backButton)
		window.SetContent(container.New(layout.NewCenterLayout(), content))
		return
	}

	form := container.New(layout.NewFormLayout())
	for name, inode := range *found {
		filePath := ""
		if dir == "/" {
			filePath = "/" + name
		} else if dir[len(dir)-1] == '/' {
			filePath = dir[:len(dir)-2] + name
		} else {
			filePath = dir + string('/') + name
		}
		foundLabel := widget.NewLabel(fmt.Sprintf("file: %v", filePath))
		restoreButton := widget.NewButton("Restore", getRestoreCallback(name, inode, worker))
		form.Add(foundLabel)
		form.Add(restoreButton)
	}
	header := widget.NewLabel(fmt.Sprintf("Found %v inodes to restore from journal", len(*found)))
	backButton := widget.NewButton("Drive list", func() { OpenSelectDrive(window) })
	content := container.New(layout.NewVBoxLayout(), header, form, backButton)
	window.SetContent(container.New(layout.NewCenterLayout(), content))
}

func OpenFindJournal(dir string, worker *ext_worker.Ext3Worker, window fyne.Window) {
	showSearch(window)
	go showFound(dir, worker, window)
}
