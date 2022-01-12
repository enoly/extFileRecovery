package gui

import (
	"fmt"
	"reflect"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	extworker "github.com/enoly/extFileRecovery/internal/extworker"
	ext3 "github.com/enoly/extFileRecovery/pkg/ext3"
)

func getExt3Worker(drive string, fs string) (*extworker.Ext3Worker, error) {
	extFs, err := ext3.New(drive)
	if err != nil {
		return nil, err
	}
	return extworker.NewExt3Worker(extFs), nil
}

func getExt3Content(drive string, fs string, window fyne.Window) *fyne.Container {
	worker, err := getExt3Worker(drive, fs)
	if err != nil {
		errorText := widget.NewLabel(fmt.Sprintf("Unable to get drive list!\n%v", err))
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
	findButton := widget.NewButton("Find indirect blocks", func() {})
	buttonsContent := container.New(layout.NewHBoxLayout(), backButton, layout.NewSpacer(), findButton)

	content := container.New(layout.NewVBoxLayout(), headerDriveName, headerDriveFs, infoBlock, buttonsContent)
	return container.New(layout.NewCenterLayout(), content)
}

func showExt3Content(drive string, fs string, window fyne.Window) {
	window.SetContent(getExt3Content(drive, fs, window))
}

func OpenExt3DriveScreen(drive string, fs string, window fyne.Window) {
	showLoading(drive, window)
	go showExt3Content(drive, fs, window)
}
