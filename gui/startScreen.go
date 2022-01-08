package gui

import "github.com/maxence-charriere/go-app/v9/pkg/app"

type StartScreen struct {
	app.Compo
}

func (s *StartScreen) Render() app.UI {
	return app.Div().Body(
		app.H1().Text("ext3-4 file recovery tool"),
		app.A().Href("/selectDrive").Body(app.Button().Body(app.Text("Recover"))),
	)
}
