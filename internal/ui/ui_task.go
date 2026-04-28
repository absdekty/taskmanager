package ui

import (
	_"context"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	_"github.com/absdekty/taskmanager/internal/service"
)

func (ui *UI) InitTask() {
	ui.sections.task = container.NewMax(ui.task.bg, widget.NewLabel("..."))
}