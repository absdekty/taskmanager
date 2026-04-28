package ui

import (
	_"context"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	_"github.com/absdekty/taskmanager/internal/service"
)

func (ui *UI) InitSearch() {
	ui.sections.search = container.NewMax(ui.search.bg, widget.NewLabel("..."))
}