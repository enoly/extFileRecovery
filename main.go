package main

import (
	"fyne.io/fyne/v2/app"
	gui "github.com/enoly/extFileRecovery/gui"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("ext3-4 file recovery tool")

	gui.OpenStartScreen(myWindow)
	myWindow.ShowAndRun()
}
