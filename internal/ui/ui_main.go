package ui

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/absdekty/taskmanager/internal/service"
)

type UI struct {
	ctx context.Context

	service service.ServiceI

	app    fyne.App
	window fyne.Window
}

func NewUI(service service.ServiceI) *UI {
	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
	w := a.NewWindow("Task manager")

	ui := &UI{
		service: service,
		app:     a,
		window:  w,
	}

	return ui
}

func (ui *UI) Run(ctx context.Context) error {
	ui.ctx = ctx

	ui.window.Resize(fyne.NewSize(1000, 800))
	ui.window.CenterOnScreen()
	ui.window.ShowAndRun()

	return nil
}
