package ui

import (
	_"context"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	_"github.com/absdekty/taskmanager/internal/service"
)

func (ui *UI) InitTags() {
	ui.sections.tags = container.NewMax(ui.tags.bg, widget.NewLabel("..."))
}