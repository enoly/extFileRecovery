package main

import (
	"log"
	"net/http"
	"os/exec"

	gui "github.com/enoly/extFileRecovery/gui"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func main() {
	app.Route("/", &gui.StartScreen{})
	app.RunWhenOnBrowser()
	http.Handle("/", &app.Handler{
		Name:        "extFileRecovery",
		Description: "ext3-4 file recovery tool",
	})

	openBrowser("http://localhost:8000/")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func openBrowser(url string) {
	if err := exec.Command("xdg-open", url).Start(); err != nil {
		log.Fatal(err)
	}
}
