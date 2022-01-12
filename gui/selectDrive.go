package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	lsblk "github.com/enoly/extFileRecovery/pkg/lsblk"
)

func getLoadingContent(header string) *fyne.Container {
	headerText := widget.NewLabel(header)
	loadingText := widget.NewLabel("Please wait...")
	headerText.Alignment = fyne.TextAlignCenter
	loadingText.Alignment = fyne.TextAlignCenter
	vbox := container.New(layout.NewVBoxLayout(), headerText, loadingText)

	return container.New(layout.NewCenterLayout(), vbox)
}

func showLoading(loadingHeader string, window fyne.Window) {
	window.SetContent(getLoadingContent(loadingHeader))
}

func getDriveList() *[]lsblk.BlockDevice {
	list, err := lsblk.GetDeviceList()
	if err != nil {
		return nil
	}

	resultList := make([]lsblk.BlockDevice, 0)
	for _, device := range list.Blockdevices {
		if device.Fstype == "ext3" || device.Fstype == "ext4" {
			resultList = append(resultList, device)
		}
	}

	return &resultList
}

func getSelectFunc(drive lsblk.BlockDevice, window fyne.Window) func() {
	driveName := drive.Name
	driveFs := drive.Fstype
	if driveFs == "ext3" {
		return func() {
			OpenExt3DriveScreen(driveName, driveFs, window)
		}
	} else {
		return func() {
			OpenExt4DriveScreen(driveName, driveFs, window)
		}
	}
}

func getDriveListContent(window fyne.Window) *fyne.Container {
	driveList := getDriveList()
	if driveList == nil {
		errorText := widget.NewLabel("Unable to get drive list!")
		errorText.Alignment = fyne.TextAlignCenter
		return container.New(layout.NewCenterLayout(), errorText)
	}

	if len(*driveList) == 0 {
		emptyText := widget.NewLabel("No ext3 or ext4 drives found")
		emptyText.Alignment = fyne.TextAlignCenter
		return container.New(layout.NewCenterLayout(), emptyText)
	}

	grid := container.New(layout.NewFormLayout())
	for _, drive := range *driveList {
		driveText := widget.NewLabel(fmt.Sprintf("%v\nfilesystem: %v", drive.Name, drive.Fstype))
		button := widget.NewButton("Select", getSelectFunc(drive, window))
		grid.Add(driveText)
		grid.Add(button)
	}

	return container.New(layout.NewCenterLayout(), grid)
}

func showDriveList(window fyne.Window) {
	window.SetContent(getDriveListContent(window))
}

func OpenSelectDrive(window fyne.Window) {
	showLoading("Select drive", window)
	go showDriveList(window)
}
