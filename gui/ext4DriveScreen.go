package gui

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	extworker "github.com/enoly/extFileRecovery/internal/extworker"
	ext4 "github.com/enoly/extFileRecovery/pkg/ext4"
)

func getExt4Worker(drive string, fs string) (*extworker.Ext4Worker, error) {
	extFs, err := ext4.New(drive)
	if err != nil {
		return nil, err
	}
	return extworker.NewExt4Worker(extFs), nil
}

func getExt4Content(drive string, fs string, window fyne.Window) *fyne.Container {
	worker, err := getExt4Worker(drive, fs)
	if err != nil {
		errorText := widget.NewLabel(fmt.Sprintf("Unable to get drive info!\n%v", err))
		errorText.Alignment = fyne.TextAlignCenter
		return container.New(layout.NewCenterLayout(), errorText)
	}

	headerDriveName := widget.NewLabel(drive)
	headerDriveName.Alignment = fyne.TextAlignCenter
	headerDriveFs := widget.NewLabel(fs)
	headerDriveFs.Alignment = fyne.TextAlignCenter

	infoBlock := container.New(layout.NewFormLayout())
	superblockInfo := reflect.ValueOf(*worker.GetSuperblockInfo())
	typeOfSB := superblockInfo.Type()
	for i := 0; i < superblockInfo.NumField(); i++ {
		fieldName := widget.NewLabel(typeOfSB.Field(i).Name)
		fieldValue := widget.NewLabel(fmt.Sprintf("%v", superblockInfo.Field(i).Interface()))
		infoBlock.Add(fieldName)
		infoBlock.Add(fieldValue)
	}

	backButton := widget.NewButton("Back", func() { OpenSelectDrive(window) })
	findButton := widget.NewButton("Find indirect blocks", func() { OpenExt4FindIndirect(worker, window) })
	buttonsContent := container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer(), findButton)

	content := container.New(layout.NewVBoxLayout(), headerDriveName, headerDriveFs, infoBlock, buttonsContent)
	return container.New(layout.NewCenterLayout(), content)
}

func showExt4Content(drive string, fs string, window fyne.Window) {
	window.SetContent(getExt4Content(drive, fs, window))
}

func OpenExt4DriveScreen(drive string, fs string, window fyne.Window) {
	showLoading(drive, window)
	go showExt4Content(drive, fs, window)
}
