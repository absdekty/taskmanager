package ui

import (
	_"context"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	_"github.com/absdekty/taskmanager/internal/service"
)

func (ui *UI) InitTask() {
	taskLabel := widget.NewLabel("Задача")

	informationContent := container.NewVBox(
		widget.NewLabel("Информация.."),
	)
	
	subtasksContent := container.NewVBox(
		widget.NewLabel("Субзадачи.."),
	)

	/* Функция обновления контента */
	ui.task.updateContent = func() {
		task := ui.currentTask
		
		if task == nil {
			taskLabel.SetText("Задача")
		} else {
			taskLabel.SetText(task.Title)
		}
	}
	
	/* Контент */
	ui.sections.task = container.NewMax(
		ui.task.bg,
		container.NewBorder(
			taskLabel,
			nil, nil, nil,
			container.NewAppTabs(
				container.NewTabItem("Информация", informationContent),
				container.NewTabItem("Субзадачи", subtasksContent))))
}