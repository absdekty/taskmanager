package ui

import (
	"context"
	"image/color"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/absdekty/taskmanager/internal/service"
	"github.com/absdekty/taskmanager/internal/model"
)

type UI struct {
	ctx context.Context

	service service.ServiceI

	app    fyne.App
	window fyne.Window
	
	currentTask *model.Task
	
	sections struct {
		search *fyne.Container
		task *fyne.Container
		tags *fyne.Container
	}
	
	search struct {
		bg *canvas.Rectangle
		
		tasks []*model.Task
		list *widget.List
		
		updateContent func()
	}
	
	task struct {
		bg *canvas.Rectangle
	}
	
	tags struct {
		bg *canvas.Rectangle
	}
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
	
	ui.InitBackground()
	
	ui.InitSearch()
	ui.InitTask()
	ui.InitTags()
	
	ui.UpdateContent()
	
	return ui
}

func (ui *UI) Run(ctx context.Context) error {
	ui.ctx = ctx
	
	ui.search.updateContent()
	
	ui.window.Resize(fyne.NewSize(1000, 800))
	ui.window.CenterOnScreen()
	ui.window.ShowAndRun()

	return nil
}

func (ui *UI) InitBackground() {
	ui.search.bg = canvas.NewRectangle(color.RGBA{43, 43, 43, 255})
	ui.task.bg = canvas.NewRectangle(color.RGBA{43, 43, 43, 255})
	ui.tags.bg = canvas.NewRectangle(color.RGBA{43, 43, 43, 255})
}

func (ui *UI) UpdateContent() {
	mainContent := container.NewHSplit(ui.sections.search, ui.sections.task)
	mainContent.SetOffset(0.25)

	ui.window.SetContent(
		container.NewBorder(
			nil,
			ui.sections.tags,
			nil, nil,
			mainContent))
}